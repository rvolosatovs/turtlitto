package main

import (
	"flag"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
	"github.com/rvolosatovs/turtlitto/pkg/webapi"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	unixSockPath = filepath.Join(os.TempDir(), "trc-sock-test")

	netLst net.Listener

	logger *zap.SugaredLogger
)

func init() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = zapLogger.Sugar()

	f, err := os.Open(unixSockPath)
	if !os.IsNotExist(err) {
		f.Close()

		logger.Infof("Removing %s...", unixSockPath)
		if err := os.Remove(unixSockPath); err != nil {
			logger.Fatalf("Failed to remove %s: %s", unixSockPath, err)
		}
	}

	netLst, err = net.Listen("unix", unixSockPath)
	if err != nil {
		logger.Fatalf("Failed to open Unix socket on %s: %s", unixSockPath, err)
	}
}

func TestMain(m *testing.M) {
	if err := flag.Set("unixSocket", unixSockPath); err != nil {
		logger.Fatalf("Failed to set `socket` to %s: %s", unixSockPath, err)
	}

	if err := flag.Set("debug", "true"); err != nil {
		logger.Fatalf("Failed to set `debug`: %s", err)
	}

	logger.Info("Starting SRRS in goroutine...")
	go main()

	dial := func() (net.Conn, error) { return net.DialTimeout("tcp", defaultAddr, time.Second) }
	retries := 20

	conn, err := dial()
	for i := 0; err != nil && i < retries; i++ {
		time.Sleep(100 * time.Millisecond)
		conn, err = dial()
	}
	if err != nil {
		logger.Fatalf("Failed to connect to SRRS at %s: %s", defaultAddr, err)
	}

	if err := conn.Close(); err != nil {
		logger.Fatalf("Failed to close connection: %s", err)
	}

	ret := m.Run()

	if err := netLst.Close(); err != nil {
		logger.Infof("Failed to close Unix socket: %s", err)
	}

	logger.Infof("Exiting with return code: %d", ret)
	os.Exit(ret)
}

func TestAll(t *testing.T) {
	a := require.New(t)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	hs := apitest.RandomHandshake()
	_ = hs

	var wsConn *websocket.Conn
	var wsErr error
	go func() {
		defer wg.Done()

		wsAddr := "ws://localhost" + defaultAddr + "/" + webapi.StateEndpoint
		logger.With("addr", wsAddr).Debug("Opening a WebSocket...")
		wsConn, _, wsErr = websocket.DefaultDialer.Dial(wsAddr, nil)
	}()

	logger.Debug("Waiting for connection on Unix socket...")
	unixConn, err := netLst.Accept()
	a.NoError(err)
	logger.Debug("Connection on Unix socket received")

	logger.Debug("Establishing mock TRC connection...")
	trc := trctest.Connect(unixConn, unixConn,
		trctest.WithHandler(api.MessageTypeHandshake, func(msg *api.Message) (*api.Message, error) {
			a.NotNil(msg.ParentID)
			a.NotEmpty(msg.ParentID)
			a.NotEmpty(msg.MessageID)
			a.NotEqual(msg.MessageID, msg.ParentID)
			a.NotEmpty(msg.Payload)
			return nil, nil
		}),
	)
	logger.Debug("Mock TRC connection established")

	logger.Debug("Sending handshake...")
	err = trc.SendHandshake(&api.Handshake{
		Version: trcapi.DefaultVersion,
	})
	a.NoError(err)
	logger.Debug("Handshake sent")

	wg.Wait()
	a.NoError(wsErr)

	err = wsConn.WriteJSON("SESSIONKEY")
	a.NoError(err)

	logger.Debug("WebSocket opened")

	var got api.State
	logger.Debug("Receiving nil state on WebSocket...")
	err = wsConn.ReadJSON(&got)
	a.Nil(err)
	a.Equal(api.State{}, got)
	logger.Debug("Nil state received on WebSocket")

	for i := 0; i < 10; i++ {
		// TODO: generate
		state := &api.State{
			Turtles: map[string]*api.TurtleState{
				"foo": {},
			},
		}

		logger.Debug("Sending random state...")
		err = trc.SendState(state)
		a.NoError(err)
		logger.Debug("Random state sent")

		got := &api.State{}
		logger.Debug("Receiving random state on WebSocket...")
		err = wsConn.ReadJSON(got)
		a.Nil(err)
		a.Equal(got, state)
		logger.Debug("Random state received on WebSocket")
	}
	// TODO: check setting of commands
}
