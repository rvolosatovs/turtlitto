package trcapi

import (
	"sync"

	"go.uber.org/zap"
)

// Pool represents a pool of Conn's.
// The Pool manages synchronisation and allows easy creation, accessing and closing of connections.
type Pool struct {
	connectFunc func() (*Conn, func(), error)
	closeFunc   func()

	connMu *sync.Mutex
	conn   *Conn
}

// NewPool returns a new Pool.
// connectFunc must return a *Conn, function to close it(possibly nil) and error, if any.
func NewPool(connectFunc func() (*Conn, func(), error)) *Pool {
	return &Pool{
		connectFunc: connectFunc,
		connMu:      &sync.Mutex{},
	}
}

// Conn returns an existing open connection, if such exists or establishes a new one.
func (p *Pool) Conn() (*Conn, error) {
	logger := zap.L()

	p.connMu.Lock()
	defer p.connMu.Unlock()

	if p.conn != nil {
		select {
		case <-p.conn.Closed():
			if p.closeFunc != nil {
				logger.Debug("Closing existing connection...")
				p.closeFunc()
			}
		default:
			logger.Debug("Returning existing open connection...")
			return p.conn, nil
		}
	}

	logger.Debug("Establishing a new open connection...")
	conn, closeFunc, err := p.connectFunc()
	if err != nil {
		return nil, err
	}
	p.conn = conn
	p.closeFunc = closeFunc
	return conn, nil
}

// Close closes the underlying connection.
func (p *Pool) Close() error {
	p.connMu.Lock()
	if p.closeFunc != nil {
		p.closeFunc()
	}
	p.conn = nil
	p.connMu.Unlock()
	return nil
}
