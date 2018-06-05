package main

import (
	"flag"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
)

var (
	debug    = flag.Bool("debug", false, "Debug mode")
	unixSock = flag.String("unixSocket", DefaultUnixSocket, "Path to the unix socket")
	tcpSock  = flag.String("tcpSocket", DefaultTCPSocket, "Service address of tcp socket. TCP will be used instead of a Unix socket when this is set")
	silent   = flag.Bool("silent", false, "Disables automatic sending of random state updates")
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
		QuoteEmptyFields: true,
	})
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	var netLst net.Listener
	var err error
	switch {
	case *unixSock != "" && *tcpSock != "":
		log.WithFields(log.Fields{
			"unixSocket": *unixSock,
			"tcpSocket":  *tcpSock,
		}).Fatal("At most one of tcpSocket and unixSocket must be specified")

	case *unixSock != "":
		logger := log.WithField("path", *unixSock)

		logger.Info("Listening on Unix socket...")
		netLst, err = net.Listen("unix", *unixSock)
		if err != nil {
			logger.WithError(err).Fatal("Failed to listen on Unix socket")
		}

	case *tcpSock != "":
		logger := log.WithField("addr", *tcpSock)

		logger.Info("Listening on TCP socket...")
		netLst, err = net.Listen("tcp", *tcpSock)
		if err != nil {
			logger.WithError(err).Fatal("Failed to listen on TCP socket")
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
						log.WithError(err).Error("Failed to close connection")
					}
				}
				return

			default:
			}

			if err != nil {
				log.WithError(err).Error("Failed to accept connection")
				continue
			}

			go func() {
				defer sockConn.Close()

				log.WithField("addr", sockConn.RemoteAddr()).Infof("Connection accepted")

				// state handler of mock TRC, logs all actions
				setStateHandler := func(msg *api.Message) (*api.Message, error) {
					log.Infof("Received: %s", msg)

					reply, err := trctest.DefaultStateHandler(msg)
					log.Debugf("Reply with: %s", reply)
					return reply, err
				}
				// ping handler of mock TRC, logs when ping is received
				pingHandler := func(msg *api.Message) (*api.Message, error) {
					log.Debug("Received ping")
					return trctest.DefaultPingHandler(msg)
				}

				trcConn := trctest.Connect(sockConn, sockConn,
					trctest.WithHandler(api.MessageTypeState, setStateHandler),
					trctest.WithHandler(api.MessageTypePing, pingHandler),
					trctest.WithHandler(api.MessageTypeHandshake, trctest.DefaultHandshakeHandler),
				)
				defer trcConn.Close()

				go func() {
					for err := range trcConn.Errors() {
						log.WithError(err).Error("Internal mock-trc error")
						return
					}
				}()

				if err := trcConn.SendHandshake(&api.Handshake{
					Version: trcapi.DefaultVersion,
					Token:   "test",
				}); err != nil {
					log.WithError(err).Error("Failed to send handshake")
					return
				}
				log.WithField("version", trcapi.DefaultVersion).Debug("Sent handshake")

				if err := trcConn.SendState(apitest.RandomState()); err != nil {
					log.WithError(err).Error("Failed to send state")
					return
				}

				if *silent {
					<-make(chan int)
					return
				}

				for {
					select {
					case <-time.After(5*time.Second + time.Millisecond*time.Duration(rand.Intn(3000))):
						if err := trcConn.SendState(apitest.RandomState()); err != nil {
							log.WithError(err).Error("Failed to send state")
							return
						}
						log.Debug("Sent state")

					case <-time.After(3*time.Second + time.Millisecond*time.Duration(rand.Intn(3000))):
						if err := trcConn.Ping(); err != nil {
							log.WithError(err).Error("Failed to send ping")
							return
						}
						log.Debug("Sent ping")

					case <-closeCh:
						return
					}
				}
			}()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c)
	sig := <-c
	close(closeCh)
	log.WithField("signal", sig).Info("Received signal, exiting...")
}
