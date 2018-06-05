package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
	"github.com/stretchr/testify/require"
)

var (
	unixSockPath = filepath.Join(os.TempDir(), "trc-sock-test")

	netLst net.Listener
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

	netLst, err = net.Listen("unix", unixSockPath)
	if err != nil {
		log.Fatalf("Failed to open Unix socket on %s: %s", unixSockPath, err)
	}
}

func TestMain(m *testing.M) {
	if err := flag.Set("unixSocket", unixSockPath); err != nil {
		log.Fatalf("Failed to set `socket` to %s: %s", unixSockPath, err)
	}

	if err := flag.Set("debug", "true"); err != nil {
		log.Fatalf("Failed to set `debug`: %s", err)
	}

	log.Print("Starting SRRS in goroutine...")
	go main()

	dial := func() (net.Conn, error) { return net.DialTimeout("tcp", defaultAddr, time.Second) }
	retries := 20

	conn, err := dial()
	for i := 0; err != nil && i < retries; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err = dial()
	}
	if err != nil {
		log.Fatalf("Failed to connect to SRRS at %s: %s", defaultAddr, err)
	}

	if err := conn.Close(); err != nil {
		log.Fatalf("Failed to close connection: %s", err)
	}

	ret := m.Run()

	if err := netLst.Close(); err != nil {
		log.Printf("Failed to close Unix socket: %s", err)
	}

	log.Printf("Exiting with return code: %d", ret)
	os.Exit(ret)
}

func TestAll(t *testing.T) {
	a := require.New(t)

	wsAddr := "ws://localhost" + defaultAddr + "/" + stateEndpoint
	log.WithField("addr", wsAddr).Debug("Opening a WebSocket...")
	wsConn, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	a.NoError(err)
	log.Debug("WebSocket opened")

	log.Debug("Waiting for connection on Unix socket...")
	unixConn, err := netLst.Accept()
	a.NoError(err)
	log.Debug("Connection on Unix socket received")

	log.Debug("Establishing mock TRC connection...")
	trc := trctest.Connect(unixConn, unixConn)
	log.Debug("Mock TRC connection established")

	log.Debug("Sending handshake...")
	err = trc.SendHandshake(&api.Handshake{
		Version: trcapi.DefaultVersion,
	})
	a.NoError(err)
	log.Debug("Handshake sent")

	var got api.State
	log.Debug("Receiving nil state on WebSocket...")
	err = wsConn.ReadJSON(&got)
	a.Nil(err)
	a.Equal(api.State{}, got)
	log.Debug("Nil state received on WebSocket")

	for i := 0; i < 10; i++ {
		// TODO: generate
		state := &api.State{
			Turtles: map[string]*api.TurtleState{
				"foo": {},
			},
		}

		log.Debug("Sending random state...")
		err = trc.SendState(state)
		a.NoError(err)
		log.Debug("Random state sent")

		got := &api.State{}
		log.Debug("Receiving random state on WebSocket...")
		err = wsConn.ReadJSON(got)
		a.Nil(err)
		a.Equal(got, state)
		log.Debug("Random state received on WebSocket")
	}
	// TODO: check setting of commands
}
