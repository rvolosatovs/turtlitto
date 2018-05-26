package main

import (
	"compress/flate"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

const (
	defaultAddr   = ":4242" // default webserver address
	retryInterval = 5 * time.Second
)

var (
	httpAddr = flag.String("http", defaultAddr, "HTTP service address")
	static   = flag.String("static", "", "Path to the static assets")
	sock     = flag.String("socket", filepath.Join(os.TempDir(), "trc.sock"), "Path to the unix socket")
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	statusEndpoint  = path.Join("api", "v1", "status")
	turtleEndpoint  = path.Join("api", "v1", "turtles")
	commandEndpoint = path.Join("api", "v1", "command")
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()

	pool := api.NewPool(func() (*api.Conn, func(), error) {
		unixConn, err := net.Dial("unix", *sock)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to connect to TRC's unix socket")
		}

		trcConn, err := api.Connect(api.DefaultVersion, unixConn, unixConn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to establish connection to TRC")
		}

		return trcConn, func() {
			if err := unixConn.Close(); err != nil {
				log.Printf("Failed to close TRC Unix socket: %s", err)
			}
			if err := trcConn.Close(); err != nil {
				log.Printf("Failed to close TRC connection: %s", err)
			}
		}, nil
	})
	defer pool.Close()

	go func() {
		var last time.Time
		for {
			time.Sleep(time.Until(last.Add(retryInterval)))
			last = time.Now()

			trcConn, err := pool.Conn()
			if err != nil {
				log.Printf("Failed to connect to TRC: %s", err)
				continue
			}

			for err := range trcConn.Errors() {
				log.Printf("Communication with TRC failed: %s", err)
			}
		}
	}()

	http.HandleFunc("/"+statusEndpoint, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// TODO: Check token

		wsConn, err := upgrader.Upgrade(w, r, nil)
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

		stateCh, closeFn := trcConn.SubscribeState()
		defer closeFn()

		for {
			select {
			case <-r.Context().Done():
				return

			case <-trcConn.Closed():
				http.Error(w, fmt.Sprintf("Communication with TRC closed"), http.StatusServiceUnavailable)
				return

			case st := <-stateCh:
				if err := wsConn.WriteJSON(st); err != nil {
					http.Error(w, fmt.Sprintf("Failed to write state: %s", err), http.StatusInternalServerError)
					return
				}
			}
		}
	})

	http.HandleFunc("/"+commandEndpoint, func(w http.ResponseWriter, r *http.Request) {
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

		if err := trcConn.SetCommand(req.Command); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/"+turtleEndpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "PUT" {
			http.Error(w, fmt.Sprintf("Expected a PUT request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		// TODO: Check token

		id := strings.TrimPrefix(r.URL.Path, "/"+turtleEndpoint+"/")
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

		if err := trcConn.SetTurtleState(id, &req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/"+turtleEndpoint, func(w http.ResponseWriter, r *http.Request) {
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
			&api.State{
				Turtles: req.Turtles,
			},
		); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusBadRequest)
			return
		}
	})

	if *static != "" {
		http.Handle("/", http.FileServer(http.Dir(*static)))
	}

	if err := (&http.Server{
		Addr:     *httpAddr,
		ErrorLog: log.New(os.Stdout, "", 0),
	}).ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen on %s: %v", *httpAddr, err)
	}
}
