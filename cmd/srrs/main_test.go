package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
	"github.com/stretchr/testify/assert"
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

	flag.Parse()
	if *tcpSock == "" {
		logger.Debug("Creating Unix socket")
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

	} else {
		logger.Debug("Creating TCP socket")
		netLst, err = net.Listen("tcp", *tcpSock)
		if err != nil {
			logger.Fatalf("Failed to open TCP socket on %s: %s", *tcpSock, err)
		}
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
	a := assert.New(t)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	hs := apitest.RandomHandshake()

	var wsConn *websocket.Conn
	var wsErr error
	go func() {
		defer wg.Done()

		wsAddr := "ws://localhost" + defaultAddr + "/" + stateEndpoint
		logger.With("addr", wsAddr).Debug("Opening a WebSocket...")
		wsConn, _, wsErr = websocket.DefaultDialer.Dial(wsAddr, http.Header{
			"Authorization": []string{fmt.Sprintf(
				"Basic %s", base64.StdEncoding.EncodeToString(append([]byte("user:"), []byte(hs.Token)...)),
			)},
		})
	}()

	logger.Debug("Waiting for connection on Unix socket...")
	unixConn, err := netLst.Accept()
	a.NoError(err)
	logger.Debug("Connection on Unix socket received")

	logger.Debug("Establishing mock TRC connection...")
	msgCh := make(chan *api.Message, 1)
	trc := trctest.Connect(
		unixConn, unixConn,
		// check for correct handshake
		trctest.WithHandler(api.MessageTypeHandshake, func(msg *api.Message) (*api.Message, error) {
			a.NotNil(msg.ParentID)
			a.NotEmpty(msg.ParentID)
			a.NotEmpty(msg.MessageID)
			a.NotEqual(msg.MessageID, msg.ParentID)
			a.NotEmpty(msg.Payload)
			return nil, nil
		}),
		// send state to channel, then react default
		trctest.WithHandler(api.MessageTypeState, func(msg *api.Message) (*api.Message, error) {
			logger.Debug("Receiving message", msg.Type)
			msgCh <- msg
			return trctest.DefaultStateHandler(msg)
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

	logger.Debug("WebSocket opened")

	var currentState api.State
	logger.Debug("Receiving nil state on WebSocket...")
	err = wsConn.ReadJSON(&currentState)
	a.Nil(err)
	a.Equal(api.State{}, currentState)
	logger.Debug("Nil state received on WebSocket")

	httpCl := http.DefaultClient
	httpCl.Timeout = 200 * time.Millisecond

	// state TRC -> SRRC
	for i := 0; i < 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			state := apitest.RandomState()

			logger.Debug("Sending random state...")
			err = trc.SendState(state)
			a.NoError(err)
			logger.Debug("Random state sent")

			got := &api.State{}
			logger.Debug("Receiving random state on WebSocket...")
			err = wsConn.ReadJSON(got)
			a.NoError(err)
			a.Equal(got, state)
			logger.Debug("Random state received on WebSocket")
		})
	}

	// state SRRC -> TRC
	for i := 0; i < 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			state := apitest.RandomState()

			logger.Debug("Sending random state SRRC -> TRC")
			pld, err := json.Marshal(state.Turtles)
			a.NoError(err)
			req, err := http.NewRequest(http.MethodPost, "http://"+defaultAddr+"/"+turtleEndpoint, bytes.NewReader(pld))
			req.Header = http.Header{
				"Authorization": []string{fmt.Sprintf(
					"Basic %s", base64.StdEncoding.EncodeToString(append([]byte("user:"), []byte(hs.Token)...)),
				)},
			}
			a.NoError(err)

			_, err = httpCl.Do(req)
			if !a.NoError(err) {
				return
			}

			got := &api.State{}
			logger.Debug("Waiting for state on TRC...")
			var msg *api.Message
			select {
			case <-time.After(httpCl.Timeout):
				t.Fail()
				return
			case msg = <-msgCh:
			}

			a.Equal(msg.Type, api.MessageTypeState)
			a.NoError(json.Unmarshal(msg.Payload, got))
			a.Equal(got, state)
			logger.Debug("Random state received on TRC")

			// wait for update from the TRC
			resp := &api.State{}
			logger.Debug("Waiting for response of TRC...")
			err = wsConn.ReadJSON(resp)
			a.NoError(err)
			a.Equal(resp, state)
			logger.Debug("Response received on WebSocket")
		})
	}

	// command TRC -> SRRC
	for i := 0; i < 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmd := apitest.RandomCommand()
			state := &api.State{
				Command: *cmd,
			}

			logger.Debugf("Sending command %s", *cmd)
			err = trc.SendState(state)
			a.NoError(err)
			logger.Debug("Random command sent")

			got := &api.State{}
			logger.Debug("Receiving random state on WebSocket...")
			err = wsConn.ReadJSON(got)
			a.NoError(err)
			a.Equal(got.Command, cmd)
			logger.Debug("Random state received on WebSocket")
		})
	}

	// command SRRC -> TRC
	for i := 0; i < 10; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmd := apitest.RandomCommand()

			logger.Debugf("Sending command %s SRRC -> TRC", *cmd)
			pld, err := json.Marshal(cmd)
			a.NoError(err)
			req, err := http.NewRequest(http.MethodPost, "http://"+defaultAddr+"/"+commandEndpoint, bytes.NewReader(pld))
			a.NoError(err)
			req.Header = http.Header{
				"Authorization": []string{fmt.Sprintf(
					"Basic %s", base64.StdEncoding.EncodeToString(append([]byte("user:"), []byte(hs.Token)...)),
				)},
			}

			_, err = httpCl.Do(req)
			if !a.NoError(err) {
				return
			}

			got := &api.State{}
			logger.Debug("Waiting for state on TRC...")
			var msg *api.Message
			select {
			case <-time.After(httpCl.Timeout):
				t.Fail()
				return
			case msg = <-msgCh:
			}

			a.Equal(msg.Type, api.MessageTypeState)
			err = json.Unmarshal(msg.Payload, got)
			a.Nil(err)
			a.Equal(got.Command, cmd)
			logger.Debug("Random state received on TRC")

			// wait for TRC response
			resp := &api.State{}
			logger.Debug("Waiting for response of TRC...")
			err = wsConn.ReadJSON(resp)
			a.NoError(err)
			a.Equal(resp, got.Command)
			logger.Debug("Response received on WebSocket")
		})
	}
}
