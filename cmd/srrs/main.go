package main

import (
	"encoding/hex"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

const defaultAddr = ":4242" // default webserver address

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

	conn, err := net.Dial("unix", *sock)
	if err != nil {
		log.Fatalf("Failed to connect to unix socket at %s: %s", *sock, err)
	}

	_ = api.NewClient(conn, conn)

	if *static != "" {
		http.Handle("/", http.FileServer(http.Dir(*static)))
	}

	http.HandleFunc("/commands", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to open WebSocket from %s: %s", r.RemoteAddr, err)
			return
		}

		for {
			mt, b, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Failed to read message from %s: %s", r.RemoteAddr, err)
				return
			}

			switch mt {
			case websocket.BinaryMessage:
				log.Printf("Got binary message from %s: %s", r.RemoteAddr, hex.EncodeToString(b))
			case websocket.TextMessage:
				log.Printf("Got text message from %s: %s", r.RemoteAddr, string(b))
			}

			if err := conn.WriteMessage(mt, []byte("pong")); err != nil {
				log.Printf("Failed to write message to %s: %s", r.RemoteAddr, err)
				return
			}
		}
	})

	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("Failed to close socket at %s: %s", *sock, err)
		}
		log.Fatalf("Failed to listen on %s: %v", *httpAddr, err)
	}
}
