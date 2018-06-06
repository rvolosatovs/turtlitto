package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/webapi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultAddr   = ":4242" // default webserver address
	retryInterval = 5 * time.Second
)

var (
	debug    = flag.Bool("debug", false, "Debug mode")
	tcpAddr  = flag.String("tcp", defaultAddr, "HTTP service address")
	static   = flag.String("static", "", "Path to the static assets")
	unixSock = flag.String("unixSocket", filepath.Join(os.TempDir(), "trc.sock"), "Path to the unix socket")
	tcpSock  = flag.String("tcpSocket", "", "Internal TCP socket address. TRC <-> SRRS communication uses TCP instead of a Unix socket when this is set")

	stateEndpoint   = path.Join("api", "v1", "state")
	turtleEndpoint  = path.Join("api", "v1", "turtles")
	commandEndpoint = path.Join("api", "v1", "command")
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

		pool := trcapi.NewPool(func() (*trcapi.Conn, func(), error) {
			var netConn net.Conn
			if *tcpSock == "" {
				logger := logger.With(zap.String("trc_socket_unix", *unixSock))

				var err error
				logger.Debug("Dialing Unix socket...")
				netConn, err = net.Dial("unix", *unixSock)
				if err != nil {
					return nil, nil, errors.Wrapf(err, "Failed to connect to TRC's unix socket")
				}
				logger.Debug("Unix socket dial succeeded")
			} else {
				logger := logger.With(zap.String("trc_socket_tcp", *tcpSock))

				var err error
				logger.Debug("Dialing TCP socket...")
				netConn, err = net.Dial("tcp", *tcpSock)
				if err != nil {
					return nil, nil, errors.Wrapf(err, "Failed to connect to TRC's TCP socket")
				}
				logger.Debug("TCP socket dial succeeded")
			}

			logger.Debug("Initializing TRC protocol connection on socket...")
			trcConn, err := trcapi.Connect(trcapi.DefaultVersion, netConn, netConn)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "Failed to establish connection to TRC")
			}
			logger.Debug("TRC protocol connection initialized")

			return trcConn, func() {
				logger.Debug("Closing TRC connection...")
				if err := trcConn.Close(); err != nil {
					logger.With(zap.Error(err)).Error("Failed to close TRC connection")
				}

				logger.Debug("Closing socket...")
				if err := netConn.Close(); err != nil {
					logger.With(zap.Error(err)).Error("Failed to close socket")
				}
			}, nil
		})
		defer pool.Close()

		go func() {
			var last time.Time
			for {
				d := time.Until(last.Add(retryInterval))
				if d > 0 {
					logger.With(zap.Duration("duration", d)).Debug("Sleeping before retrying the connection...")
				}
				time.Sleep(d)

				last = time.Now()

				logger.Debug("Retrieving a connection from pool...")
				trcConn, err := pool.Conn()
				if err != nil {
					logger.With(zap.Error(err)).Warn("Failed to retrieve TRC connection from pool")
					continue
				}
				logger.Debug("Connection retrieval from pool succeeded")

				for err := range trcConn.Errors() {
					logger.With(zap.Error(err)).Error("Communication with TRC failed")
				}
			}
		}()

		http.HandleFunc("/"+stateEndpoint, webapi.MakeStateHandler(pool))
		http.HandleFunc("/"+commandEndpoint, webapi.MakeCommandHandler(pool))
		http.HandleFunc("/"+turtleEndpoint+"/", webapi.MakeTurtleHandler(pool))

		if *static != "" {
			http.Handle("/", http.FileServer(http.Dir(*static)))
		}

		logger := logger.With(zap.String("listen_addr_tcp", *tcpAddr))

		logger.Info("Starting the web server...")
		if err := (&http.Server{
			Addr:     *tcpAddr,
			ErrorLog: zap.NewStdLog(logger),
			Handler: &webapi.LogHandler{
				Logger:  logger,
				Handler: http.DefaultServeMux,
			},
		}).ListenAndServe(); err != nil {
			return errors.Wrap(err, "failed to listen")
		}
		return nil

	}(); err != nil {
		logger.With(zap.Error(err)).Fatal("SRRS failed")
	}
}
