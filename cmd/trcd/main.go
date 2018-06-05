package main

import (
	"flag"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
)

var (
	debug    = flag.Bool("debug", false, "Debug mode")
	unixSock = flag.String("unixSocket", filepath.Join(os.TempDir(), "trc.sock"), "Path to the unix socket")
	tcpSock  = flag.String("tcpSocket", "", "Service address of tcp socket. TCP will be used instead of a Unix socket when this is set")
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
	if *tcpSock == "" {
		log.Info("Listening to Unix socket...")
		netLst, err = net.Listen("unix", *unixSock)
	} else {
		log.Info("Listening to TCP socket...")
		addr, _ := net.ResolveTCPAddr("", *tcpSock)
		netLst, err = net.ListenTCP("tcp", addr)
	}

	if err != nil {
		panic(err)
	}
	defer netLst.Close()

	for {
		sockConn, err := netLst.Accept()
		log.Infof("Connected to %s", sockConn.RemoteAddr())
		if err != nil {
			panic(err)
		}
		defer sockConn.Close()

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

		go func() {
			for err := range trcConn.Errors() {
				log.WithError(err).Error("Internal mock-trc error")
			}
		}()

		if err := trcConn.SendHandshake(&api.Handshake{Version: trcapi.DefaultVersion}); err != nil {
			panic(err)
		}
		log.Debug("Received handshake")

		if err := trcConn.SendState(&api.State{}); err != nil {
			log.Error("Initial state: " + err.Error())
			return
		}

		closeCh := make(chan struct{})
		defer close(closeCh)

		// send random messages unless silent is set
		if !(*silent) {
			go func() {
				hBeat := time.Tick(3 * time.Second)
				for {
					aWhile := time.Millisecond * time.Duration(rand.Intn(5000))
					// wait 0-5 seconds, then send an empty state
					select {
					case <-time.After(aWhile):
						if err := trcConn.SendState(apitest.RandomState()); err != nil {
							log.Error("Random sendState error: " + err.Error())
						} else {
							log.Debug("Sent random state")
						}
					case <-hBeat:
						if err := trcConn.Ping(); err != nil {
							log.Error("Ping error: " + err.Error())
						} else {
							log.Debug("Sent ping")
						}
					case <-closeCh:
						return
					}
				}
			}()
		}
	}
}
