package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/stretchr/testify/assert"
)

var (
	unixSockPath = filepath.Join(os.TempDir(), "trc-sock-test")

	unixSock net.Listener
)

func init() {
	f, err := os.Open(unixSockPath)
	if !os.IsNotExist(err) {
		f.Close()

		log.Printf("Removing %s...", unixSockPath)
		if err := os.Remove(unixSockPath); err != nil {
			log.Fatalf("Failed to remove %s: %s", unixSockPath, err)
		}
	}

	unixSock, err = net.Listen("unix", unixSockPath)
	if err != nil {
		log.Fatalf("Failed to open Unix socket on %s: %s", unixSockPath, err)
	}
}

func TestMain(m *testing.M) {
	if err := flag.Set("socket", unixSockPath); err != nil {
		log.Fatalf("Failed to set `socket` to %s: %s", unixSockPath, err)
	}
	go func() {
		log.Print("Starting SRRS in goroutine...")
		main()
	}()

	ret := m.Run()

	if err := unixSock.Close(); err != nil {
		log.Printf("Failed to close Unix socket: %s", err)
	}

	log.Printf("Exiting with return code: %d", ret)
	os.Exit(ret)
}

func TestAll(t *testing.T) {
	a := assert.New(t)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost"+defaultAddr+"/state", nil)
	if !a.Nil(err) {
		t.FailNow()
	}

	state := map[string]*api.State{
		"foo": &api.State{
			ID: "bar",
			// TODO: add more fields
		},
	}

	go func() {
		for {
			c, err := unixSock.Accept()
			a.Nil(err)

			var req api.Message

			err = json.NewDecoder(c).Decode(&req)
			if !a.Nil(err) {
				return
			}

			resp := &api.Message{
				MessageID: req.MessageID,
				Type:      req.Type,
			}

			switch req.Type {
			case api.MessageTypeSetState:
			case api.MessageTypeGetState:
				b, err := json.Marshal(state)
				if err != nil {
					panic(err)
				}
				resp.Payload = b
			case api.MessageTypeCommand:
			default:
				t.Errorf("Unmatched message type: %s", req.Type)
				return
			}

			err = json.NewEncoder(c).Encode(resp)
			a.Nil(err)
		}
	}()

	var got map[string]*api.State

	err = conn.ReadJSON(&got)
	if !a.Nil(err) {
		t.FailNow()
	}
	a.Equal(state, got)
}
