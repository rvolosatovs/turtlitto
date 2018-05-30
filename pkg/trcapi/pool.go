package trcapi

import "sync"

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
	p.connMu.Lock()
	defer p.connMu.Unlock()

	if p.conn != nil {
		select {
		case <-p.conn.Closed():
			if p.closeFunc != nil {
				p.closeFunc()
			}
		default:
			return p.conn, nil
		}
	}

	conn, closeFunc, err := p.connectFunc()
	if err != nil {
		return nil, err
	}
	p.conn = conn
	p.closeFunc = closeFunc
	return conn, nil
}

func (p *Pool) Close() error {
	p.closeFunc()
	return nil
}
