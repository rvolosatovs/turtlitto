package trcapi

import (
	"sync"

	"go.uber.org/zap"
)

type Pool struct {
	connectFunc func() (*Conn, func(), error)
	closeFunc   func()

	connMu *sync.Mutex
	conn   *Conn
}

func NewPool(connectFunc func() (*Conn, func(), error)) *Pool {
	return &Pool{
		connectFunc: connectFunc,
		connMu:      &sync.Mutex{},
	}
}

func (p *Pool) Conn() (*Conn, error) {
	logger := zap.L().With(zap.String("func", "trcapi.Pool.Conn"))

	logger.Debug("Locking connMu...")
	p.connMu.Lock()
	defer p.connMu.Unlock()

	logger.Debug("connMu locked...")

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

func (p *Pool) Close() error {
	p.connMu.Lock()
	if p.closeFunc != nil {
		p.closeFunc()
	}
	p.connMu.Unlock()
	return nil
}
