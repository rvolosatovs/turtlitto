package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/webapi"
)

const (
	defaultAddr   = ":4242" // default webserver address
	retryInterval = 5 * time.Second
)

var (
	httpAddr = flag.String("http", defaultAddr, "HTTP service address")
	static   = flag.String("static", "", "Path to the static assets")
	sock     = flag.String("socket", filepath.Join(os.TempDir(), "trc.sock"), "Path to the unix socket")

	statusEndpoint  = path.Join("api", "v1", "status")
	turtleEndpoint  = path.Join("api", "v1", "turtles")
	commandEndpoint = path.Join("api", "v1", "command")
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()

	pool := trcapi.NewPool(func() (*trcapi.Conn, func(), error) {
		unixConn, err := net.Dial("unix", *sock)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to connect to TRC's unix socket")
		}

		trcConn, err := trcapi.Connect(trcapi.DefaultVersion, unixConn, unixConn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to establish connection to TRC")
		}

		return trcConn, func() {
			if err := unixConn.Close(); err != nil {
				log.Printf("Failed to close TRC Unix socket: %s", err)
			}
			if err := trcConn.Close(); err != nil {
				log.Printf("Failed to close TRC connection: %s", err)
			}
		}, nil
	})
	defer pool.Close()

	go func() {
		var last time.Time
		for {
			time.Sleep(time.Until(last.Add(retryInterval)))
			last = time.Now()

			trcConn, err := pool.Conn()
			if err != nil {
				log.Printf("Failed to connect to TRC: %s", err)
				continue
			}

			for err := range trcConn.Errors() {
				log.Printf("Communication with TRC failed: %s", err)
			}
		}
	}()

	http.HandleFunc("/"+statusEndpoint, webapi.StateHandler(pool))
	http.HandleFunc("/"+commandEndpoint, webapi.CommandHandler(pool))
	http.HandleFunc("/"+turtleEndpoint+"/", webapi.TurtlesHandler(pool))
	http.HandleFunc("/"+turtleEndpoint, webapi.TurtlesIDHandler(pool, func(r *http.Request) string {
		return strings.TrimPrefix(r.URL.Path, "/"+turtleEndpoint+"/")
	}))

	if *static != "" {
		http.Handle("/", http.FileServer(http.Dir(*static)))
	}

	if err := (&http.Server{
		Addr:     *httpAddr,
		ErrorLog: log.New(os.Stdout, "", 0),
	}).ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen on %s: %v", *httpAddr, err)
	}
}
