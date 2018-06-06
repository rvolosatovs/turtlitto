package main

import (
	"flag"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	debug    = flag.Bool("debug", false, "Debug mode")
	unixSock = flag.String("unixSocket", DefaultUnixSocket, "Path to the unix socket")
	tcpSock  = flag.String("tcpSocket", DefaultTCPSocket, "Service address of tcp socket. TCP will be used instead of a Unix socket when this is set")
	silent   = flag.Bool("silent", false, "Disables automatic sending of random state updates")
)

func main() {
	flag.Parse()

	var logger *zap.Logger
	var err error
	if *debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "logger",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}), os.Stdout, zap.DebugLevel)).Fatal("Failed to initialize logger")
	}
	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)

	if err := func() error {
		defer logger.Sync() //nolint

		var netLst net.Listener
		switch {
		case *unixSock != "" && *tcpSock != "":
			return errors.New("At most one of tcpSocket and unixSocket must be specified")

		case *unixSock != "":
			logger := logger.With(zap.String("path", *unixSock))

			logger.Info("Listening on Unix socket...")
			netLst, err = net.Listen("unix", *unixSock)
			if err != nil {
				return errors.Wrap(err, "failed to listen on Unix socket")
			}

		case *tcpSock != "":
			logger := logger.With(zap.String("addr", *tcpSock))

			logger.Info("Listening on TCP socket...")
			netLst, err = net.Listen("tcp", *tcpSock)
			if err != nil {
				return errors.Wrap(err, "failed to listen on TCP socket")
			}
		}

		defer netLst.Close()

		closeCh := make(chan struct{})

		go func() {
			for {
				sockConn, err := netLst.Accept()

				select {
				case <-closeCh:
					if sockConn != nil {
						err = sockConn.Close()
						if err != nil {
							logger.With(zap.Error(err)).Error("Failed to close connection")
						}
					}
					return

				default:
				}

				if err != nil {
					logger.With(zap.Error(err)).Error("Failed to accept connection")
					continue
				}

				go func() {
					defer sockConn.Close()

					logger := logger.With(zap.Stringer("addr", sockConn.RemoteAddr()))

					logger.Info("Connection accepted")

					trcConn := trctest.Connect(sockConn, sockConn,
						trctest.WithHandler(api.MessageTypeState, func(msg *api.Message) (*api.Message, error) {
							logger.With(zap.Any("state", msg)).Info("Received state")

							reply, err := trctest.DefaultStateHandler(msg)
							logger.With(zap.Any("reply", reply)).Debug("Sending reply...")
							return reply, err
						}),

						trctest.WithHandler(api.MessageTypePing, func(msg *api.Message) (*api.Message, error) {
							logger.Debug("Received ping")
							return trctest.DefaultPingHandler(msg)
						}),

						trctest.WithHandler(api.MessageTypeHandshake, func(msg *api.Message) (*api.Message, error) {
							logger.Debug("Received handshake")
							return trctest.DefaultPingHandler(msg)
						}),
					)
					defer trcConn.Close()

					go func() {
						for err := range trcConn.Errors() {
							logger.With(zap.Error(err)).Error("Internal mock-trc error")
							return
						}
					}()

					if err := trcConn.SendHandshake(&api.Handshake{
						Version: trcapi.DefaultVersion,
						Token:   "test",
					}); err != nil {
						logger.With(zap.Error(err)).Error("Failed to send handshake")
						return
					}
					logger.With(zap.Stringer("version", trcapi.DefaultVersion)).Debug("Sent handshake")

					if err := trcConn.SendState(apitest.RandomState()); err != nil {
						logger.With(zap.Error(err)).Error("Failed to send state")
						return
					}
					logger.Info("Sent initial state")

					if *silent {
						<-make(chan int)
						return
					}

					go func() {
						for {
							select {
							case <-time.After(5*time.Second + time.Millisecond*time.Duration(rand.Intn(3000))):
								if err := trcConn.SendState(apitest.RandomState()); err != nil {
									logger.With(zap.Error(err)).Error("Failed to send state")
									return
								}
								logger.Debug("Sent state")

							case <-closeCh:
								return
							}
						}
					}()

					go func() {
						for {
							select {
							case <-time.After(3*time.Second + time.Millisecond*time.Duration(rand.Intn(3000))):
								if err := trcConn.Ping(); err != nil {
									logger.With(zap.Error(err)).Error("Failed to send ping")
									continue
								}
								logger.Debug("Sent ping")

							case <-closeCh:
								return
							}
						}
					}()
				}()
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		sig := <-c
		close(closeCh)
		logger.With(zap.Stringer("signal", sig)).Info("Received signal, exiting...")
		return nil
	}(); err != nil {
		logger.With(zap.Error(err)).Fatal("TRCD failed")
	}
}
