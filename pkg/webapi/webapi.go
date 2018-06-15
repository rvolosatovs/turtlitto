// Package webapi implements the web API as defined in the specification.
package webapi

import (
	"compress/flate"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/logcontext"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"go.uber.org/zap"
)

var (
	// StateEndpoint is the state endpoint.
	StateEndpoint = path.Join("api", "v1", "state")

	// AuthEndpoint is the authentication endpoint.
	AuthEndpoint = path.Join("api", "v1", "auth")

	// TurtleEndpoint is the turtle endpoint.
	TurtleEndpoint = path.Join("api", "v1", "turtles")

	// CommandEndpoint is the command endpoint.
	CommandEndpoint = path.Join("api", "v1", "command")

	errActiveWebSocket     = errors.New("an active WebSocket connection already exists")
	errAuthenticateFirst   = errors.New("authenticate first")
	errAuthorizationHeader = errors.New("Authorization header not found or invalid")
	errInvalidSessionKey   = errors.New("invalid session key")
	errInvalidToken        = errors.New("invalid token")
	errFailedToGetToken    = errors.New("TRC connection established, but failed to get token")
)

// controlWriter can write Control messages to itself.
type controlWriter interface {
	WriteControl(messageType int, data []byte, deadline time.Time) error
}

// wsError closes websocket represented by w with and error message err and code code.
// wsError logs to logger.
func wsError(w controlWriter, logger *zap.Logger, err error, code int) {
	logger = logger.With(
		zap.Error(err),
		zap.Int("code", code),
	)

	logger.Debug("Closing WebSocket...")
	if err := w.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, err.Error()), time.Time{}); err != nil {
		logger.Error("Failed to gracefully close WebSocket")
	}
}

type session struct {
	isActive bool
	key      string
}

// server manages the web API.
type server struct {
	pool *trcapi.Pool

	sessionMu sync.RWMutex
	session   *session
}

// handleState handles requests to StateEndpoint.
func (srv *server) handleState(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logcontext.Logger(ctx)

	wsConn, err := (&websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(w, r, nil)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to open WebSocket")
		return
	}
	defer wsConn.Close()

	wsConn.EnableWriteCompression(true)
	if err := wsConn.SetCompressionLevel(flate.BestCompression); err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to enable compression"), websocket.CloseProtocolError)
		return
	}

	var key string
	logger.Debug("Reading key...")

	if err := wsConn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to set read deadline"), websocket.CloseInternalServerErr)
		return
	}

	if err := wsConn.ReadJSON(&key); err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to read session key"), websocket.CloseInvalidFramePayloadData)
		return
	}

	srv.sessionMu.Lock()
	switch {
	case srv.session == nil:
		wsError(wsConn, logger, errAuthenticateFirst, websocket.ClosePolicyViolation)
		srv.sessionMu.Unlock()
		return

	case key != srv.session.key:
		wsError(wsConn, logger, errInvalidSessionKey, websocket.CloseInvalidFramePayloadData)
		srv.sessionMu.Unlock()
		return

	case srv.session.isActive:
		wsError(wsConn, logger, errActiveWebSocket, websocket.ClosePolicyViolation)
		srv.sessionMu.Unlock()
		return
	}

	srv.session.isActive = true
	srv.sessionMu.Unlock()

	defer func() {
		srv.sessionMu.Lock()
		srv.session.isActive = false
		srv.sessionMu.Unlock()
	}()

	logger.Debug("Retrieving a connection from pool...")
	trcConn, err := srv.pool.Conn()
	if err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to establish connection to TRC"), websocket.CloseInternalServerErr)
		return
	}

	logger.Debug("Subscribing to state changes...")
	changeCh, closeFn, err := trcConn.SubscribeStateChanges(ctx)
	if err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to subscribe to state changes"), websocket.CloseInternalServerErr)
		return
	}
	defer closeFn()

	oldState := trcConn.State(ctx)

	logger.Debug("Sending current state on the WebSocket...")
	if err := wsConn.WriteJSON(oldState); err != nil {
		logger.With(zap.Error(err)).Error("Failed to write state")
		return
	}
	for {
		select {
		case <-ctx.Done():
			wsError(wsConn, logger, errors.New("Context done"), websocket.CloseInvalidFramePayloadData)
			return

		case <-trcConn.Closed():
			wsError(wsConn, logger, errors.New("TRC connection is closed"), websocket.CloseInternalServerErr)
			return

		case <-changeCh:
			logger.Debug("State change acknowledged")

			st := trcConn.State(ctx)
			// TODO: Compute diff of st and oldState
			_ = oldState
			diff := st
			oldState = st
			logger.Debug("Sending state diff on the WebSocket...")
			if err := wsConn.WriteJSON(diff); err != nil {
				wsError(wsConn, logger, errors.Wrap(err, "failed to write state"), websocket.CloseInternalServerErr)
				return
			}
			logger.Debug("Sending state diff on the WebSocket succeeded")
		}
	}
}

// handleAuth handles requests to AuthEndpoint.
func (srv *server) handleAuth(w http.ResponseWriter, r *http.Request) {
	logger := logcontext.Logger(r.Context())

	if r.Method != "GET" {
		http.Error(w, errors.Errorf("expected a GET request, got %s", r.Method).Error(), http.StatusBadRequest)
		return
	}

	srv.sessionMu.Lock()
	defer srv.sessionMu.Unlock()

	if srv.session != nil && srv.session.isActive {
		http.Error(w, errActiveWebSocket.Error(), http.StatusTeapot)
		return
	}

	logger.Debug("Retrieving a connection from pool...")
	trcConn, err := srv.pool.Conn()
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to establish connection to TRC").Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("Retrieving token...")
	trcTok, err := trcConn.Token()
	if err != nil {
		http.Error(w, errors.Wrap(err, errFailedToGetToken.Error()).Error(), http.StatusInternalServerError)
		return
	}

	_, authTok, ok := r.BasicAuth()
	if !ok && trcTok != "" {
		http.Error(w, errAuthorizationHeader.Error(), http.StatusBadRequest)
		return
	}

	if trcTok != "" && authTok != trcTok {
		http.Error(w, errInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	logger.Debug("Generating new session key...")
	b := make([]byte, 64)
	_, err = rand.Read(b)
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to generate session key").Error(), http.StatusInternalServerError)
		return
	}

	key := hex.EncodeToString(b)
	_, err = w.Write([]byte(key))
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to write session key").Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("Creating new session")
	srv.session = &session{
		key: key,
	}
}

func (srv *server) makeTRCSendHandler(f func(context.Context, *trcapi.Conn, *json.Decoder) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logcontext.Logger(ctx)

		var err error
		if r.Method != "POST" {
			http.Error(w, errors.Errorf("Expected a POST request, got %s", r.Method).Error(), http.StatusBadRequest)
			return
		}

		srv.sessionMu.RLock()
		defer srv.sessionMu.RUnlock()

		_, key, ok := r.BasicAuth()
		if !ok && srv.session.key != "" {
			http.Error(w, errAuthorizationHeader.Error(), http.StatusBadRequest)
			return
		}

		switch {
		case srv.session == nil:
			http.Error(w, errAuthenticateFirst.Error(), http.StatusMethodNotAllowed)
			return

		case key != srv.session.key:
			http.Error(w, errInvalidSessionKey.Error(), http.StatusUnauthorized)
			return
		}

		logger.Debug("Retrieving a connection from pool...")
		trcConn, err := srv.pool.Conn()
		if err != nil {
			http.Error(w, errors.Wrap(err, "failed to establish connection to TRC").Error(), http.StatusInternalServerError)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := f(ctx, trcConn, dec); err != nil {
			http.Error(w, errors.Wrap(err, "failed to process request").Error(), http.StatusBadRequest)
			return
		}
	}
}

// HandleFuncer allows registration of a handler function for a specified pattern.
// An example implementation of this interface is *http.ServeMux.
type HandleFuncer interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// Register endpoints registers webapi endpoints on handler.
func RegisterHandlers(pool *trcapi.Pool, handler HandleFuncer) {
	s := &server{
		pool: pool,
	}

	for ep, f := range map[string]http.HandlerFunc{
		"/" + AuthEndpoint: s.handleAuth,

		"/" + StateEndpoint: s.handleState,

		"/" + CommandEndpoint: s.makeTRCSendHandler(func(ctx context.Context, trcConn *trcapi.Conn, dec *json.Decoder) error {
			var cmd api.Command
			if err := dec.Decode(&cmd); err != nil {
				return errors.Wrap(err, "failed to decode request body")
			}

			if err := trcConn.SetCommand(ctx, cmd); err != nil {
				return errors.Wrap(err, "failed to send command to TRC")
			}
			return nil
		}),

		"/" + TurtleEndpoint: s.makeTRCSendHandler(func(ctx context.Context, trcConn *trcapi.Conn, dec *json.Decoder) error {
			var st map[string]*api.TurtleState
			if err := dec.Decode(&st); err != nil {
				return errors.Wrap(err, "failed to read states")
			}

			if err := trcConn.SetTurtleState(ctx, st); err != nil {
				fmt.Println("TURTLE", err)
				return errors.Wrap(err, "failed to send turtle state to TRC")
			}
			return nil
		}),
	} {
		handler.HandleFunc(ep, f)
	}
}
