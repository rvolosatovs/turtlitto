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
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"
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

	logger *zap.Logger
)

const (
	timeout      = time.Second
	messageCount = 3
)

func init() {
	http.DefaultClient.Timeout = timeout

	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	flag.Parse()
	if *tcpSock == "" {
		sockLogger := logger.With(zap.String("path", unixSockPath))
		sockLogger.Debug("Creating Unix socket...")
		f, err := os.Open(unixSockPath)
		if !os.IsNotExist(err) {
			f.Close()

			sockLogger.Info("Removing UNIX socket...")
			if err := os.Remove(unixSockPath); err != nil {
				sockLogger.Fatal("Failed to remove UNIX socket",
					zap.Error(err),
				)
			}
		}

		sockLogger.Debug("Listening on UNIX socket...")
		netLst, err = net.Listen("unix", unixSockPath)
		if err != nil {
			sockLogger.Fatal("Failed to listen on UNIX socket",
				zap.Error(err),
			)
		}

	} else {
		sockLogger := logger.With(zap.String("addr", *tcpSock))

		sockLogger.Debug("Listening on TCP socket...")
		netLst, err = net.Listen("tcp", *tcpSock)
		if err != nil {
			sockLogger.Fatal("Failed to listen on TCP socket",
				zap.Error(err),
			)
		}
	}
}

func TestMain(m *testing.M) {
	if err := flag.Set("unixSocket", unixSockPath); err != nil {
		logger.Fatal("Failed to set `socket` to "+unixSockPath,
			zap.Error(err),
		)
	}

	if err := flag.Set("debug", "true"); err != nil {
		logger.Fatal("Failed to set `debug`",
			zap.Error(err),
		)
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
		logger.Fatal("Failed to connect to SRRS",
			zap.String("addr", defaultAddr),
			zap.Error(err),
		)
	}

	if err := conn.Close(); err != nil {
		logger.Fatal("Failed to close connection",
			zap.Error(err),
		)
	}

	logger = zap.L()

	ret := m.Run()

	if err := netLst.Close(); err != nil {
		logger.Info("Failed to close Unix socket", zap.Error(err))
	}

	logger.Info("Exiting with return code", zap.Int("code", ret))
	os.Exit(ret)
}

func TestAPI(t *testing.T) {
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
			logger.Error("Failed to authenticate",
				zap.Error(err),
			)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		a.NoError(err)

		sessionKey = string(b)
		logger.Debug("Got session key",
			zap.String("key", sessionKey),
		)
	}()

	logger.Debug("Waiting for connection on Unix socket...")
	unixConn, err := netLst.Accept()
	a.NoError(err)
	defer unixConn.Close()

	logger.Debug("Establishing mock TRC connection...")
	msgCh := make(chan *api.Message)
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
	logger.Debug("Opening a WebSocket...",
		zap.String("addr", wsAddr),
	)
	wsConn, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if !a.NoError(err) {
		t.Fatal("Failed to open WebSocket")
	}
	defer wsConn.Close()

	err = wsConn.WriteJSON(sessionKey)
	a.NoError(err)

	currentState := &api.State{}

	if !t.Run("TRC->SRRC/state", func(t *testing.T) {
		for i := 0; i < messageCount; i++ {
			if !t.Run(strconv.Itoa(i), func(t *testing.T) {
				a = assert.New(t)

				st := apitest.RandomState()
				if err := st.Validate(); err != nil {
					panic(errors.Wrap(err, "invalid state generated"))
				}

				logger.Debug("Sending random state from TRC...",
					zap.Reflect("state", st),
				)
				err = trc.SendState(st)
				a.NoError(err)

				b, err := json.Marshal(st)
				if err != nil {
					panic(err)
				}

				merged := deepcopy.Copy(currentState).(*api.State)
				if err = json.Unmarshal(b, merged); err != nil {
					panic(err)
				}

				st, err = api.StateDiff(currentState, merged)
				if err != nil {
					panic(err)
				}
				if st == nil {
					st = &api.State{}
				}

				got := &api.State{}
				err = wsConn.ReadJSON(got)
				a.NoError(err)
				a.Equal(st, got)
			}) {
				t.FailNow()
			}
		}
	}) {
		t.FailNow()
	}

	if !t.Run("SRRC->TRC/turtles", func(t *testing.T) {
		for i := 0; i < messageCount; i++ {
			if !t.Run(strconv.Itoa(i), func(t *testing.T) {
				a = assert.New(t)

				st := &api.State{
					Turtles: apitest.RandomTurtleStateMap(),
				}
				for len(st.Turtles) == 0 {
					st.Turtles = apitest.RandomTurtleStateMap()
				}
				if err := st.Validate(); err != nil {
					panic(errors.Wrap(err, "invalid state generated"))
				}

				b, err := json.Marshal(st.Turtles)
				a.NoError(err)

				req, err := http.NewRequest(http.MethodPost, "http://"+defaultAddr+"/"+webapi.TurtleEndpoint, bytes.NewReader(b))
				a.NoError(err)
				req.SetBasicAuth("", sessionKey)

				errCh := make(chan error, 1)
				go func() {
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						errCh <- err
					}
					defer resp.Body.Close()

					a.Equal(resp.StatusCode, http.StatusOK)

					b, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						errCh <- errors.Wrap(err, "failed to read response body")
					}

					if len(b) > 0 {
						errCh <- errors.Errorf("server returned error: %s", string(b))
					}
					errCh <- nil
				}()

				var msg *api.Message
				select {
				case <-time.After(timeout):
					t.Fatal("Timed out waiting for message to arrive at SRRS")
				case msg = <-msgCh:
				}

				a.Equal(msg.Type, api.MessageTypeState)
				a.Nil(msg.ParentID)
				a.NotEmpty(msg.MessageID)

				select {
				case <-time.After(timeout):
					t.Fatal("Timed out sending state to SRRS")
				case err = <-errCh:
					if !a.NoError(err) {
						return
					}
				}

				merged := deepcopy.Copy(currentState).(*api.State)
				if err = json.Unmarshal(msg.Payload, merged); err != nil {
					panic(err)
				}

				st, err = api.StateDiff(currentState, merged)
				if err != nil {
					panic(err)
				}
				if st == nil {
					st = &api.State{}
				}

				got := &api.State{}
				err = wsConn.ReadJSON(got)
				a.NoError(err)
				a.Equal(st, got)

				wg.Wait()
			}) {
				t.FailNow()
			}
		}
	}) {
		t.FailNow()
	}

	t.Run("SRRC->TRC/command", func(t *testing.T) {
		for i := 0; i < messageCount; i++ {
			if !t.Run(strconv.Itoa(i), func(t *testing.T) {
				a = assert.New(t)

				st := &api.State{
					Command: apitest.RandomCommand(),
				}

				b, err := json.Marshal(st.Command)
				a.NoError(err)

				req, err := http.NewRequest(http.MethodPost, "http://"+defaultAddr+"/"+webapi.CommandEndpoint, bytes.NewReader(b))
				a.NoError(err)
				req.SetBasicAuth("", sessionKey)

				errCh := make(chan error, 1)
				go func() {
					logger.Debug("Sending command to SRRS...")
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						errCh <- err
					}
					defer resp.Body.Close()

					a.Equal(resp.StatusCode, http.StatusOK)

					b, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						errCh <- errors.Wrap(err, "failed to read response body")
					}

					if len(b) > 0 {
						errCh <- errors.Errorf("server returned error: %s", string(b))
					}
					errCh <- nil
				}()

				var msg *api.Message
				select {
				case <-time.After(timeout):
					t.Fatal("Timed out waiting for message to arrive at SRRS")
				case msg = <-msgCh:
				}

				a.Equal(msg.Type, api.MessageTypeState)
				a.Nil(msg.ParentID)
				a.NotEmpty(msg.MessageID)

				select {
				case <-time.After(timeout):
					t.Fatal("Timed out sending command to SRRS")
				case err = <-errCh:
					if !a.NoError(err) {
						return
					}
				}

				merged := deepcopy.Copy(currentState).(*api.State)
				if err = json.Unmarshal(msg.Payload, merged); err != nil {
					panic(err)
				}

				st, err = api.StateDiff(currentState, merged)
				if err != nil {
					panic(err)
				}
				if st == nil {
					st = &api.State{}
				}

				got := &api.State{}
				err = wsConn.ReadJSON(got)
				a.NoError(err)
				a.Equal(st, got)

				wg.Wait()
			}) {
				t.FailNow()
			}
		}
	})
}
