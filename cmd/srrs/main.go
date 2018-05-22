package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

const (
	defaultAddr     = ":4242" // default webserver address
	defaultEndpoint = "trc"
	maxTurtles      = 10
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
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()

	http.HandleFunc("/"+defaultEndpoint, func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open WebSocket: %s", err), http.StatusBadRequest)
			return
		}
		defer wsConn.Close()

		unixConn, err := net.Dial("unix", *sock)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to connect to TRC's unix socket: %s", err), http.StatusInternalServerError)
			return
		}
		defer unixConn.Close()

		cl := api.NewClient(unixConn, unixConn)

		go func() {
			for {
				select {
				case <-r.Context().Done():
					return
				default:
				}

				var cmd api.Command
				if err := wsConn.ReadJSON(&cmd); err != nil {
					http.Error(w, fmt.Sprintf("Failed to read command: %s", err), http.StatusBadRequest)
					continue
				}

				if err = cl.SendCommand(cmd); err != nil {
					http.Error(w, fmt.Sprintf("Failed to send command to TRC: %s", err), http.StatusInternalServerError)
					return
				}
			}
		}()

		for range time.Tick(time.Second) {
			select {
			case <-r.Context().Done():
				return
			default:
			}

			st, err := cl.State()
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get state from TRC: %s", err), http.StatusInternalServerError)
				return
			}

			if err := wsConn.WriteJSON(st); err != nil {
				http.Error(w, fmt.Sprintf("Failed to write state: %s", err), http.StatusInternalServerError)
				return
			}
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
