package main

import (
	"os"
	"path/filepath"
)

// +build darwin linux

var (
	DefaultUnixSocket = filepath.Join(os.TempDir(), "trc.sock")
	DefaultTCPSocket  = ""
)
