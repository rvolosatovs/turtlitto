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
	"github.com/rvolosatovs/turtlitto"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

const (
	defaultAddr = ":4242" // default webserver address
	maxTurtles  = 10
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

	unixConn, err := net.Dial("unix", *sock)
	if err != nil {
		log.Fatalf("Failed to connect to unix socket at %s: %s", *sock, err)
	}
	defer unixConn.Close()

	cl := api.NewClient(unixConn, unixConn)

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to open WebSocket: %s", err), http.StatusBadRequest)
			return
		}
		defer wsConn.Close()

		go func() {
			for {
				select {
				case <-r.Context().Done():
					return
				default:
				}

				st := make(map[string]*turtlitto.State, maxTurtles)
				if err := wsConn.ReadJSON(&st); err != nil {
					http.Error(w, fmt.Sprintf("Failed to read state: %s", err), http.StatusBadRequest)
					continue
				}

				st, err = cl.SetState(st)
				if err != nil {
					http.Error(w, fmt.Sprintf("Failed to send state to TRC: %s", err), http.StatusInternalServerError)
					continue
				}

				if err := wsConn.WriteJSON(st); err != nil {
					http.Error(w, fmt.Sprintf("Failed to write state: %s", err), http.StatusInternalServerError)
					fmt.Println(err)
					continue
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
				continue
			}

			if err := wsConn.WriteJSON(st); err != nil {
				http.Error(w, fmt.Sprintf("Failed to write state: %s", err), http.StatusInternalServerError)
				fmt.Println(err)
				continue
			}
		}
	})

	if *static != "" {
		http.Handle("/", http.FileServer(http.Dir(*static)))
	}

	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		if cerr := unixConn.Close(); cerr != nil {
			log.Printf("Failed to close socket at %s: %s", *sock, err)
		}
		log.Fatalf("Failed to listen on %s: %v", *httpAddr, err)
	}
}
