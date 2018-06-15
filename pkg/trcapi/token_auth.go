// +build !noauth

// Package trcapi implements the TRC API as defined in the specification.
package trcapi

import (
	"github.com/pkg/errors"
)

// Token returns the token received from TRC during the handshake procedure or error,
// if it did not happen yet.
func (c *Conn) Token() (string, error) {
	v := c.token.Load()
	if v == nil {
		return "", errors.New("No token configured")
	}

	tok, ok := v.(string)
	if !ok {
		panic("Token is not a string")
	}
	return tok, nil
}
