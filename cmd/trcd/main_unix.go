// +build darwin linux

package main

import (
	"os"
	"path/filepath"
)

var (
	DefaultUnixSocket = filepath.Join(os.TempDir(), "trc.sock")
	DefaultTCPSocket  = ""
)
