package webapi

import (
	"compress/flate"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
)

func StateHandler(pool *trcapi.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// TODO: Check token

		wsConn, err := (&websocket.Upgrader{
			EnableCompression: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}).Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open WebSocket: %s", err), http.StatusBadRequest)
			return
		}
		defer wsConn.Close()

		wsConn.EnableWriteCompression(true)
		wsConn.SetCompressionLevel(flate.BestCompression)

		trcConn, err := pool.Conn()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to establish connection to TRC: %s", err), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		changeCh, closeFn, err := trcConn.SubscribeStateChanges(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to subscribe to state changes: %s", err), http.StatusInternalServerError)
			return
		}
		defer closeFn()

		oldState := trcConn.State(ctx)

		log.Debug("Sending current state on the WebSocket...")
		if err := wsConn.WriteJSON(oldState); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write state: %s", err), http.StatusInternalServerError)
			return
		}
		for {
			select {
			case <-ctx.Done():
				log.Debug("Context done")
				return

			case <-trcConn.Closed():
				http.Error(w, fmt.Sprintf("Communication with TRC closed"), http.StatusServiceUnavailable)
				return

			case <-changeCh:
				log.Debug("State change acknowledged")

				st := trcConn.State(ctx)
				// TODO: Compute diff of st and oldState
				_ = oldState
				diff := st
				oldState = st
				log.Debug("Sending state diff on the WebSocket...")
				if err := wsConn.WriteJSON(diff); err != nil {
					http.Error(w, fmt.Sprintf("Failed to write state: %s", err), http.StatusInternalServerError)
					return
				}
				log.Debug("Sending state diff on the WebSocket succeeded")
			}
		}
	}
}

func CommandHandler(pool *trcapi.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "PUT" {
			http.Error(w, fmt.Sprintf("Expected a PUT request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		// TODO: Check token

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var req struct {
			Command api.Command
		}
		if err := dec.Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to read command: %s", err), http.StatusBadRequest)
			return
		}

		trcConn, err := pool.Conn()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to establish connection to TRC: %s", err), http.StatusInternalServerError)
			return
		}

		if err := trcConn.SetCommand(r.Context(), req.Command); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	}
}

func TurtlesHandler(pool *trcapi.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "PUT" {
			http.Error(w, fmt.Sprintf("Expected a PUT request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		// TODO: Check token

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		defer r.Body.Close()

		var req struct {
			Turtles map[string]*api.TurtleState
		}
		if err := dec.Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to read states: %s", err), http.StatusBadRequest)
			return
		}

		trcConn, err := pool.Conn()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to establish connection to TRC: %s", err), http.StatusInternalServerError)
			return
		}

		if err := trcConn.SetState(
			r.Context(),
			&api.State{
				Turtles: req.Turtles,
			},
		); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	}
}

func TurtlesIDHandler(pool *trcapi.Pool, parseID func(r *http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "PUT" {
			http.Error(w, fmt.Sprintf("Expected a PUT request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		// TODO: Check token
		id := parseID(r)
		if id == "" {
			http.Error(w, fmt.Sprintf("An ID must be specified"), http.StatusBadRequest)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var req api.TurtleState
		if err := dec.Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to read state: %s", err), http.StatusBadRequest)
			return
		}

		trcConn, err := pool.Conn()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to establish connection to TRC: %s", err), http.StatusInternalServerError)
			return
		}

		if err := trcConn.SetTurtleState(r.Context(), id, &req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	}
}

type Option func(*Server)

type Server struct {
}
