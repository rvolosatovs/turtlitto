package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
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

	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if *debug {
		conf = zap.NewDevelopmentConfig()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := conf.Build()
	if err != nil {
		panic(err)
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
					logger.Error("Failed to accept connection",
						zap.Error(err),
					)
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
							logger.Error("Internal TRCD error",
								zap.Error(err),
							)
							return
						}
					}()

					hs := &api.Handshake{
						Version: trcapi.DefaultVersion,
						Token:   "test",
					}
					if err := trcConn.SendHandshake(hs); err != nil {
						logger.Error("Failed to send handshake",
							zap.Error(err),
						)
						return
					}
					logger.Info("Sent handshake",
						zap.Reflect("handshake", hs),
					)

					st := apitest.RandomState()
					if err := trcConn.SendState(st); err != nil {
						logger.Error("Failed to send initial state",
							zap.Error(err),
						)
						return
					}
					logger.Info("Sent initial state",
						zap.Reflect("state", st),
					)

					if *silent {
						<-closeCh
						return
					}

					wg := &sync.WaitGroup{}
					wg.Add(2)

					go func() {
						defer wg.Done()

						for {
							select {
							case <-time.After(10*time.Second + time.Millisecond*time.Duration(rand.Intn(7000))):
								st := apitest.RandomState()
								if err := trcConn.SendState(st); err != nil {
									logger.Error("Failed to send state",
										zap.Error(err),
									)
									return
								}
								logger.Info("Sent state",
									zap.Reflect("state", st),
								)

							case <-closeCh:
								logger.Debug("TRCD closed, stopping state-sending goroutine")
								return
							}
						}
					}()

					go func() {
						defer wg.Done()

						for {
							select {
							case <-time.After(time.Second + time.Millisecond*time.Duration(rand.Intn(3000))):
								if err := trcConn.Ping(); err != nil {
									logger.Error("Failed to send ping",
										zap.Error(err),
									)
									return
								}
								logger.Info("Sent ping")

							case <-closeCh:
								logger.Debug("TRCD closed, stopping ping-sending goroutine")
								return
							}
						}
					}()

					wg.Wait()
				}()
			}
		}()

		fmt.Println(`********************************************************************************
                                        TOKEN INCOMING...
********************************************************************************`)
		time.Sleep(time.Second)
		fmt.Println(`********************************************************************************
********************************************************************************
                              STILL COMING...
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************`)
		time.Sleep(time.Second)
		fmt.Println(`********************************************************************************
********************************************************************************
                              STILL COMING...
********************************************************************************`)
		time.Sleep(time.Second)
		fmt.Println(`********************************************************************************
****************BEHOLD........UNLEASHING THE TOKEN.....*************************
********************************************************************************
********************************************************************************
********************************************************************************
********************************************************************************`)
		time.Sleep(time.Second)
		fmt.Println(`********************************************************************************
********************************************************************************
                  TOKEN IS:                test
********************************************************************************
********************************************************************************
********************************************************************************
		`)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		select {
		case <-closeCh:
		case sig := <-c:
			close(closeCh)
			logger.Info("Received signal, exiting...",
				zap.Stringer("signal", sig),
			)
		}
		return nil
	}(); err != nil {
		logger.With(zap.Error(err)).Fatal("TRCD failed")
	}
}
