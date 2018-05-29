package main

import (
	"encoding/json"
	"flag"
	"io"
	"net"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

var (
	sock = flag.String("socket", filepath.Join(os.TempDir(), "trc.sock"), "Path to the unix socket")
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

	conn, err := net.Listen("unix", *sock)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		c, err := conn.Accept()
		if err != nil {
			panic(err)
		}
		defer c.Close()

		go func() {
			dec := json.NewDecoder(c)
			dec.DisallowUnknownFields()
			var msg api.Message
			if err := dec.Decode(&msg); err != nil {
				log.Printf("Failed to decode input: %s", err)
			}
			log.Printf("Got: %v", msg)
		}()
		if _, err := io.Copy(c, os.Stdin); err != nil {
			log.Printf("Failed to copy input: %s", err)
			return
		}
	}
}
