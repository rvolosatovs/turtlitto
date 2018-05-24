package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

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

		trcConn, err := api.Connect(api.DefaultVersion, unixConn, unixConn)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to establish connection to TRC: %s", err), http.StatusInternalServerError)
			return
		}
		defer trcConn.Close()

		stateCh, closeFn := trcConn.SubscribeState()
		defer closeFn()

		for {
			select {
			case <-r.Context().Done():
				return
			case st := <-stateCh:
				if err := wsConn.WriteJSON(st); err != nil {
					http.Error(w, fmt.Sprintf("Failed to write state: %s", err), http.StatusInternalServerError)
					return
				}
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
