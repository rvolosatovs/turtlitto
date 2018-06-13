package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
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
	"github.com/rvolosatovs/turtlitto/pkg/webapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	unixSockPath = filepath.Join(os.TempDir(), "trc-sock-test")

	netLst net.Listener

	logger *zap.SugaredLogger
)

const (
	timeout      = time.Second
	messageCount = 3
)

func init() {
	http.DefaultClient.Timeout = timeout

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = zapLogger.Sugar()

	flag.Parse()
	if *tcpSock == "" {
		logger.Debug("Creating Unix socket...")
		f, err := os.Open(unixSockPath)
		if !os.IsNotExist(err) {
			f.Close()

			logger.Infof("Removing %s...", unixSockPath)
			if err := os.Remove(unixSockPath); err != nil {
				logger.Fatalf("Failed to remove %s: %s", unixSockPath, err)
			}
		}

		logger.Debugf("Listening on UNIX socket on %s...", unixSockPath)
		netLst, err = net.Listen("unix", unixSockPath)
		if err != nil {
			logger.Fatalf("Failed to open UNIX socket on %s: %s", unixSockPath, err)
		}

	} else {
		logger.Debugf("Listening on TCP socket on %s...", *tcpSock)
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

	logger = zap.S()

	ret := m.Run()

	if err := netLst.Close(); err != nil {
		logger.Infof("Failed to close Unix socket: %s", err)
	}

	logger.Infof("Exiting with return code: %d", ret)
	os.Exit(ret)
}

func TestAll(t *testing.T) {
	a := assert.New(t)

	handshake := apitest.RandomHandshake()

	var sessionKey string

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		req, err := http.NewRequest(http.MethodGet, "http://"+defaultAddr+"/"+webapi.AuthEndpoint, nil)
		a.NoError(err)
		req.SetBasicAuth("", handshake.Token)

		logger.Debug("Sending authentication request...")
		resp, err := http.DefaultClient.Do(req)
		if !a.NoError(err) {
			logger.With("error", err).Error("Failed to authenticate")
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		a.NoError(err)

		logger.With("key", string(b)).Debug("Got session key")
		sessionKey = string(b)
	}()

	logger.Debug("Waiting for connection on Unix socket...")
	unixConn, err := netLst.Accept()
	a.NoError(err)
	defer unixConn.Close()

	logger.Debug("Establishing mock TRC connection...")
	msgCh := make(chan *api.Message, 1)
	trc := trctest.Connect(
		unixConn, unixConn,
		trctest.WithHandler(api.MessageTypePing, func(msg *api.Message) (*api.Message, error) {
			return trctest.DefaultPingHandler(msg)
		}),
		trctest.WithHandler(api.MessageTypeHandshake, func(msg *api.Message) (*api.Message, error) {
			a.NotNil(msg.ParentID)
			a.NotEmpty(msg.ParentID)
			a.NotEmpty(msg.MessageID)
			a.NotEqual(msg.MessageID, msg.ParentID)
			a.NotEmpty(msg.Payload)
			return nil, nil
		}),
		trctest.WithHandler(api.MessageTypeState, func(msg *api.Message) (*api.Message, error) {
			msgCh <- msg
			return trctest.DefaultStateHandler(msg)
		}),
	)

	logger.Debug("Sending handshake...")
	err = trc.SendHandshake(&api.Handshake{
		Version: trcapi.DefaultVersion,
	})
	a.NoError(err)

	logger.Debug("Waiting for authentication...")
	wg.Wait()
	if t.Failed() {
		t.FailNow()
	}

	wsAddr := "ws://localhost" + defaultAddr + "/" + webapi.StateEndpoint
	logger.With("addr", wsAddr).Debug("Opening a WebSocket...")
	wsConn, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if !a.NoError(err) {
		t.FailNow()
	}
	defer wsConn.Close()

	err = wsConn.WriteJSON(sessionKey)
	a.NoError(err)

	t.Run("TRC->SRRC/state", func(t *testing.T) {
		for i := 0; i < messageCount; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				expected := apitest.RandomState()

				logger.Debug("Sending random state from TRC...")
				err = trc.SendState(expected)
				a.NoError(err)

				var got api.State
				logger.Debug("Receiving random state on WebSocket...")
				err = wsConn.ReadJSON(&got)
				a.NoError(err)
				a.Equal(expected, &got)
			})
		}
	})

	t.Run("SRRC->TRC/turtles", func(t *testing.T) {
		for i := 0; i < messageCount; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				expected := &api.State{
					Turtles: apitest.RandomState().Turtles,
				}

				b, err := json.Marshal(expected.Turtles)
				a.NoError(err)

				req, err := http.NewRequest(http.MethodPost, "http://"+defaultAddr+"/"+webapi.TurtleEndpoint, bytes.NewReader(b))
				a.NoError(err)
				req.SetBasicAuth("", sessionKey)

				wg := &sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer wg.Done()

					logger.Debug("Sending state to SRRS...")
					_, err = http.DefaultClient.Do(req)
					if !a.NoError(err) {
						return
					}
				}()

				var msg *api.Message
				select {
				case <-time.After(timeout):
					t.FailNow()
				case msg = <-msgCh:
				}

				a.Equal(msg.Type, api.MessageTypeState)
				a.Nil(msg.ParentID)
				a.NotEmpty(msg.MessageID)

				var got api.State
				err = json.Unmarshal(msg.Payload, &got)
				a.NoError(err)
				a.Equal(expected, &got)

				got = api.State{}
				err = wsConn.ReadJSON(&got)
				a.NoError(err)
				a.Equal(expected, &got)

				wg.Wait()
			})
		}
	})

	t.Run("SRRC->TRC/command", func(t *testing.T) {
		for i := 0; i < messageCount; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				expected := &api.State{
					Command: *apitest.RandomCommand(),
				}

				b, err := json.Marshal(expected.Command)
				a.NoError(err)

				req, err := http.NewRequest(http.MethodPost, "http://"+defaultAddr+"/"+webapi.CommandEndpoint, bytes.NewReader(b))
				a.NoError(err)
				req.SetBasicAuth("", sessionKey)

				wg := &sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer wg.Done()

					logger.Debug("Sending command to SRRS...")
					_, err = http.DefaultClient.Do(req)
					if !a.NoError(err) {
						return
					}
				}()

				var msg *api.Message
				select {
				case <-time.After(timeout):
					t.FailNow()
				case msg = <-msgCh:
				}

				a.Equal(msg.Type, api.MessageTypeState)
				a.Nil(msg.ParentID)
				a.NotEmpty(msg.MessageID)

				var got api.State
				err = json.Unmarshal(msg.Payload, &got)
				a.NoError(err)
				a.Equal(expected, &got)

				got = api.State{}
				err = wsConn.ReadJSON(&got)
				a.NoError(err)
				a.Equal(expected, &got)

				wg.Wait()
			})
		}
	})
}
