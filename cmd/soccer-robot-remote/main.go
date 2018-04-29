package main

import (
	"encoding/hex"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const defaultAddr = ":4242" // default webserver address

var (
	httpAddr = flag.String("http", defaultAddr, "HTTP service address")
	static   = flag.String("static", "", "Path to the static assets")
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	flag.Parse()

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
		log.Fatalf("Failed to listen on %s: %v", *httpAddr, err)
	}
}
