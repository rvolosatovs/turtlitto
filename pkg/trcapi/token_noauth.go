// +build noauth

// Package trcapi implements the TRC API as defined in the specification.
package trcapi

import "go.uber.org/zap"

// Token returns the token received from TRC during the handshake procedure or error,
// if it did not happen yet.
func (c *Conn) Token() (string, error) {
	zap.L().Warn("Bypassing TRC authentication")
	return "", nil
}
