package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

const (
	defaultAddr     = ":42424" // default webserver address
	defaultEndpoint = "relay"
	defaultCommand  = "cd /srr && docker-compose pull && docker-compose up -d"
)

var (
	httpAddr = flag.String("http", defaultAddr, "HTTP service address")
	endpoint = flag.String("endpoint", defaultEndpoint, "Endpoint to listen for requests on")
	script   = flag.String("script", defaultCommand, "Command to run once a matching request is received on the endpoint")
	token    = flag.String("token", "", "Command to run once a matching request is received on the endpoint")
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()

	http.HandleFunc("/"+*endpoint, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "GET" {
			http.Error(w, fmt.Sprintf("Expected GET, got %s", r.Method), http.StatusMethodNotAllowed)
			return
		}

		t := r.Header.Get("token")
		switch {
		case t == "":
			http.Error(w, fmt.Sprintf("Empty token"), http.StatusUnauthorized)
			return
		case t != *token:
			http.Error(w, fmt.Sprintf("Invalid token"), http.StatusUnauthorized)
			return
		}

		out, err := exec.Command("/usr/bin/env", "bash", "-c", *script).CombinedOutput()
		if err != nil {
			log.Printf("Command failed: %s", err)
			http.Error(w, fmt.Sprintf("Failed to run command: %s", err), http.StatusInternalServerError)
			return
		}
		log.Printf("Command finished with success. Output: %s", out)

		_, err = w.Write(out)
		if err != nil {
			log.Printf("Failed to write output: %s", err)
		}
	})

	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Failed to listen on %s: %v", *httpAddr, err)
	}
}
