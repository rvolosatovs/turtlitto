// Package webapi implements the web API as defined in the specification.
package webapi

import (
	"compress/flate"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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

const (
	// pingInterval is pingInterval.
	pingInterval = 5 * time.Second

	// writeTimeout is writeTimeout.
	writeTimeout = 3 * time.Second

	// readTimeout is readTimeout.
	readTimeout = 3 * time.Second

	// inactivityTimeout is inactivityTimeout.
	inactivityTimeout = 5 * time.Second
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
	errAuthorizationHeader = errors.New("`Authorization` header not found or invalid")
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

	logger.Error("Closing WebSocket...")
	if err := w.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, err.Error()), time.Now().Add(writeTimeout)); err != nil {
		logger.Warn("Failed to gracefully close WebSocket")
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

	stopTimerMu sync.Mutex
	stopTimer   *time.Timer
}

// handleState handles requests to StateEndpoint.
func (srv *server) handleState(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logcontext.Logger(ctx)

	wsConn, err := (&websocket.Upgrader{
		HandshakeTimeout:  readTimeout,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Failed to open WebSocket")
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

	if err := wsConn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
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
		srv.sessionMu.Unlock()
		wsError(wsConn, logger, errAuthenticateFirst, websocket.ClosePolicyViolation)
		return

	case key != srv.session.key:
		srv.sessionMu.Unlock()
		wsError(wsConn, logger, errInvalidSessionKey, websocket.CloseInvalidFramePayloadData)
		return

	case srv.session.isActive:
		srv.sessionMu.Unlock()
		wsError(wsConn, logger, errActiveWebSocket, websocket.ClosePolicyViolation)
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

	if err := wsConn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to set write deadline"), websocket.CloseInternalServerErr)
		return
	}

	logger.Debug("Sending current state on the WebSocket...", zap.Reflect("state", oldState))
	if err := wsConn.WriteJSON(oldState); err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to write state"), websocket.CloseInternalServerErr)
		return
	}

	if err := wsConn.SetReadDeadline(time.Now().Add(pingInterval + writeTimeout + readTimeout)); err != nil {
		wsError(wsConn, logger, errors.Wrap(err, "failed to set read deadline"), websocket.CloseInternalServerErr)
	}
	wsConn.SetPongHandler(func(string) error {
		return wsConn.SetReadDeadline(time.Now().Add(pingInterval + writeTimeout + readTimeout))
	})

	errCh := make(chan error, 1)
	go func() {
		for {
			_, _, err := wsConn.NextReader()
			if err != nil {
				errCh <- err
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wsError(wsConn, logger, errors.New("context done"), websocket.CloseInvalidFramePayloadData)
			return

		case <-trcConn.Closed():
			wsError(wsConn, logger, errors.New("TRC connection is closed"), websocket.CloseInternalServerErr)
			return

		case <-trcConn.Errors():
			defer trcConn.Close()
			wsError(wsConn, logger, errors.New("communication with TRC failed"), websocket.CloseInternalServerErr)
			return

		case err := <-errCh:
			wsError(wsConn, logger, errors.Wrap(err, "communication via WebSocket failed"), websocket.CloseAbnormalClosure)
			return

		case <-changeCh:
			logger.Debug("State change acknowledged")

			st := trcConn.State(ctx)
			// TODO: Compute diff of st and oldState
			_ = oldState

			diff := st
			if diff == nil {
				continue
			}
			oldState = st

			if err := wsConn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
				wsError(wsConn, logger, errors.Wrap(err, "failed to set write deadline"), websocket.CloseInternalServerErr)
				return
			}

			logger.Debug("Sending state diff on the WebSocket...", zap.Reflect("state", diff))
			if err := wsConn.WriteJSON(diff); err != nil {
				wsError(wsConn, logger, errors.Wrap(err, "failed to write state"), websocket.CloseInternalServerErr)
				return
			}

		case <-time.After(pingInterval):
			if err := wsConn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
				wsError(wsConn, logger, errors.Wrap(err, "failed to set write deadline"), websocket.CloseInternalServerErr)
				return
			}

			if err := wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				wsError(wsConn, logger, errors.Wrap(err, "failed to write ping"), websocket.CloseInternalServerErr)
				return
			}
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
	var stopTimerMu sync.Mutex
	stopTimer := time.AfterFunc(420 /* blaze it */, func() {
		trcConn, err := pool.Conn()
		if err != nil {
			zap.L().Error("Failed to establish connection to TRC", zap.Error(err))
			return
		}
		defer trcConn.Close()

		if err := trcConn.SetCommand(context.Background(), api.CommandStop); err != nil {
			zap.L().Error("Failed to stop TRC", zap.Error(err))
		}
	})
	stopTimer.Stop()

	activeConns := 0

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
			if cmd == "" {
				return nil
			}

			zap.L().Info("Received command", zap.String("command", string(cmd)))
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
			if len(st) == 0 {
				return nil
			}

			zap.L().Info("Received turtle state", zap.Reflect("state", st))
			if err := trcConn.SetTurtleState(ctx, st); err != nil {
				return errors.Wrap(err, "failed to send turtle state to TRC")
			}
			return nil
		}),
	} {
		hdl := f
		handler.HandleFunc(ep, func(w http.ResponseWriter, r *http.Request) {
			stopTimerMu.Lock()
			activeConns++
			stopTimer.Stop()
			stopTimerMu.Unlock()

			hdl(w, r)

			stopTimerMu.Lock()
			activeConns--
			if activeConns == 0 {
				stopTimer.Reset(inactivityTimeout)
			}
			stopTimerMu.Unlock()
		})
	}
}
