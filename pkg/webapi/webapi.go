package webapi

import (
	"compress/flate"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"go.uber.org/zap"
)

type logKey struct{}

type logResponseWriter struct {
	http.ResponseWriter
	statusCh   chan int
	responseCh chan []byte
}

func (w *logResponseWriter) Write(p []byte) (int, error) {
	b := make([]byte, len(p))
	copy(b, p)
	w.responseCh <- b
	return w.ResponseWriter.Write(p)
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.statusCh <- status
	w.ResponseWriter.WriteHeader(status)
}

type LogHandler struct {
	http.Handler
	*zap.Logger
}

type hijackResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
}

func (h LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger.With(
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("uri", r.RequestURI),
	)

	responseCh := make(chan []byte, 1)
	statusCh := make(chan int, 1)

	var wrapped http.ResponseWriter = &logResponseWriter{
		ResponseWriter: w,
		responseCh:     responseCh,
		statusCh:       statusCh,
	}
	if h, ok := w.(http.Hijacker); ok {
		wrapped = &hijackResponseWriter{
			ResponseWriter: wrapped,
			Hijacker:       h,
		}
	}

	logger.Debug("Processing request...")
	h.Handler.ServeHTTP(wrapped, r.WithContext(
		context.WithValue(r.Context(), logKey{}, logger),
	))

	status := <-statusCh
	logger = logger.With(
		zap.String("response", string(<-responseCh)),
		zap.Int("status", status),
	)
	if status < http.StatusBadRequest {
		logger.Debug("Successfully processed request")
	} else {
		logger.Error("Error processing request")
	}
}

type ControlWriter interface {
	WriteControl(messageType int, data []byte, deadline time.Time) error
}

func wsError(w ControlWriter, logger *zap.Logger, err error, code int) {
	logger = logger.With(
		zap.Error(err),
		zap.Int("code", code),
	)

	logger.Debug("Closing WebSocket...")
	if err := w.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, err.Error()), time.Time{}); err != nil {
		logger.Error("Failed to gracefully close WebSocket")
	}
}

func StateHandler(pool *trcapi.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx := r.Context()

		logger, ok := ctx.Value(logKey{}).(*zap.Logger)
		if !ok || logger == nil {
			logger = zap.L()
		}

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
		wsConn.SetCompressionLevel(flate.BestCompression)

		logger.Debug("Retrieving a connection from pool...")
		trcConn, err := pool.Conn()
		if err != nil {
			wsError(wsConn, logger, errors.Wrap(err, "failed to establish connection to TRC"), websocket.CloseInternalServerErr)
			return
		}

		logger.Debug("Retrieving token...")
		tok, err := trcConn.Token()
		if err != nil {
			wsError(wsConn, logger, errors.Wrap(err, "TRC connection established, but failed to get token"), websocket.CloseInternalServerErr)
			return
		}

		_, pswd, ok := r.BasicAuth()
		if !ok {
			wsError(wsConn, logger, errors.New("Authorization header not found or invalid"), websocket.CloseInvalidFramePayloadData)
			return
		}

		if tok != "" && pswd != tok {
			wsError(wsConn, logger, errors.New("invalid token"), websocket.CloseInvalidFramePayloadData)
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
}

func CommandHandler(pool *trcapi.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx := r.Context()

		logger, ok := ctx.Value(logKey{}).(*zap.Logger)
		if !ok || logger == nil {
			logger = zap.L()
		}

		if r.Method != "PUT" {
			http.Error(w, fmt.Sprintf("Expected a PUT request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var cmd api.Command
		if err := dec.Decode(&cmd); err != nil {
			http.Error(w, fmt.Sprintf("Failed to read command: %s", err), http.StatusBadRequest)
			return
		}

		trcConn, err := pool.Conn()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to establish connection to TRC: %s", err), http.StatusInternalServerError)
			return
		}

		tok, err := trcConn.Token()
		if err != nil {
			logger.With(zap.Error(err)).Warn("TRC connection established, but failed to get token")
			http.Error(w, fmt.Sprintf("Token not initialized"), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("token") != tok {
			logger.With(zap.Error(err)).Debug("Authentication failed")
			http.Error(w, fmt.Sprintf("Wrong authentication token"), http.StatusUnauthorized)
			return
		}

		if err := trcConn.SetCommand(r.Context(), cmd); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	}
}

func TurtleHandler(pool *trcapi.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx := r.Context()

		logger, ok := ctx.Value(logKey{}).(*zap.Logger)
		if !ok || logger == nil {
			logger = zap.L()
		}

		if r.Method != "PUT" {
			http.Error(w, fmt.Sprintf("Expected a PUT request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		defer r.Body.Close()

		var st map[string]*api.TurtleState
		if err := dec.Decode(&st); err != nil {
			http.Error(w, fmt.Sprintf("Failed to read states: %s", err), http.StatusBadRequest)
			return
		}

		trcConn, err := pool.Conn()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to establish connection to TRC: %s", err), http.StatusInternalServerError)
			return
		}

		tok, err := trcConn.Token()
		if err != nil {
			logger.With(zap.Error(err)).Warn("TRC connection established, but failed to get token")
			http.Error(w, fmt.Sprintf("Token not initialized"), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("token") != tok {
			logger.With(zap.Error(err)).Debug("Authentication failed")
			http.Error(w, fmt.Sprintf("Wrong authentication token"), http.StatusUnauthorized)
			return
		}

		if err := trcConn.SetTurtleState(r.Context(), st); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	}
}
