package main

import (
	"flag"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/webapi"
)

const (
	defaultAddr   = ":4242" // default webserver address
	retryInterval = 5 * time.Second
)

var (
	debug    = flag.Bool("debug", false, "Debug mode")
	httpAddr = flag.String("http", defaultAddr, "HTTP service address")
	static   = flag.String("static", "", "Path to the static assets")
	sock     = flag.String("socket", filepath.Join(os.TempDir(), "trc.sock"), "Path to the unix socket")

	stateEndpoint   = path.Join("api", "v1", "state")
	turtleEndpoint  = path.Join("api", "v1", "turtles")
	commandEndpoint = path.Join("api", "v1", "command")
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

	pool := trcapi.NewPool(func() (*trcapi.Conn, func(), error) {
		logger := log.WithField("path", *sock)

		logger.Debug("Dialing Unix socket...")
		unixConn, err := net.Dial("unix", *sock)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to connect to TRC's unix socket")
		}
		logger.Debug("Unix socket dial succeeded")

		logger.Debug("Initializing TRC protocol connection on Unix socket...")
		trcConn, err := trcapi.Connect(trcapi.DefaultVersion, unixConn, unixConn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to establish connection to TRC")
		}
		logger.Debug("TRC protocol connection initialized")

		return trcConn, func() {
			logger.WithError(err).Debug("Closing TRC connection...")
			if err := trcConn.Close(); err != nil {
				logger.WithError(err).Error("Failed to close TRC connection")
			}

			logger.WithError(err).Debug("Closing Unix socket...")
			if err := unixConn.Close(); err != nil {
				logger.WithError(err).Error("Failed to close Unix socket")
			}
		}, nil
	})
	defer pool.Close()

	go func() {
		var last time.Time
		for {
			d := time.Until(last.Add(retryInterval))
			if d > 0 {
				log.WithField("duration", d).Debug("Sleeping before retrying the connection...")
			}
			time.Sleep(d)

			last = time.Now()

			log.Debug("Retrieving a connection from pool...")
			trcConn, err := pool.Conn()
			if err != nil {
				log.WithError(err).Print("Failed to retrieve TRC connection from pool")
				continue
			}
			log.Debug("Connection retrieval from pool succeeded")

			for err := range trcConn.Errors() {
				log.WithError(err).Print("Communication with TRC failed")
			}
		}
	}()

	http.HandleFunc("/"+stateEndpoint, webapi.StateHandler(pool))
	http.HandleFunc("/"+commandEndpoint, webapi.CommandHandler(pool))
	http.HandleFunc("/"+turtleEndpoint+"/", webapi.TurtlesHandler(pool))
	http.HandleFunc("/"+turtleEndpoint, webapi.TurtlesIDHandler(pool, func(r *http.Request) string {
		return strings.TrimPrefix(r.URL.Path, "/"+turtleEndpoint+"/")
	}))

	if *static != "" {
		http.Handle("/", http.FileServer(http.Dir(*static)))
	}

	logger := log.WithField("addr", *httpAddr)

	w := logger.Writer()
	defer w.Close()

	logger.Info("Starting the web server...")
	if err := (&http.Server{
		Addr:     *httpAddr,
		ErrorLog: stdlog.New(w, "", 0),
	}).ListenAndServe(); err != nil {
		logger.WithError(err).Fatal("Failed to listen")
	}
}
