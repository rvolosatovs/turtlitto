package trctest

import (
	"encoding/json"
	"io"

	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

// Default handler for SetState messages, replying according to the API.
func DefaultStateHandler(msg *api.Message) (*api.Message, error) {
	var st api.State
	if err := json.Unmarshal(msg.Payload, &st); err != nil {
		return nil, err
	}

	b, err := json.Marshal(&api.State{})
	if err != nil {
		return nil, err
	}
	return api.NewMessage(api.MessageTypeState, b, &msg.MessageID), nil
}

// Default handler for ping messages, replying according to the API.
func DefaultPingHandler(msg *api.Message) (*api.Message, error) {
	return api.NewMessage(api.MessageTypePing, nil, &msg.MessageID), nil
}

// Handler is a function, which handles a message.
type Handler func(*api.Message) (*api.Message, error)

type Conn struct {
	decoder interface{ Decode(v interface{}) error }
	encoder interface{ Encode(v interface{}) error }

	errCh   chan error
	closeCh chan struct{}

	handlers      map[api.MessageType]Handler
	defaultHander Handler
}

type Option func(*Conn)

func WithHandler(msg api.MessageType, h Handler) Option {
	return func(c *Conn) {
		c.handlers[msg] = h
	}
}

func WithDefaultHandler(h Handler) Option {
	return func(c *Conn) {
		c.defaultHander = h
	}
}

// Connect establishes the TRC-side connection according to TRC API protocol
// specification of version ver on w and r.
func Connect(w io.Writer, r io.Reader, opts ...Option) (*Conn, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	conn := &Conn{
		decoder: dec,
		encoder: json.NewEncoder(w),
		closeCh: make(chan struct{}),
		errCh:   make(chan error),
	}
	for _, opt := range opts {
		opt(conn)
	}

	if conn.defaultHander == nil {
		conn.defaultHander = func(msg *api.Message) (*api.Message, error) {
			return nil, errors.Errorf("Unmatched handler for type %s", msg.Type)
		}
	}

	go func() {
		for {
			var msg api.Message
			err := conn.decoder.Decode(&msg)

			select {
			case <-conn.closeCh:
				// Don't handle err if connection is closed
				close(conn.errCh)
				return
			}
			if err != nil {
				conn.errCh <- errors.Wrap(err, "Failed to decode incoming message")
				return
			}

			h, ok := conn.handlers[msg.Type]
			if !ok {
				h = conn.defaultHander
			}

			resp, err := h(&msg)
			if err != nil {
				conn.errCh <- err
				return
			}

			if err := conn.encoder.Encode(resp); err != nil {
				conn.errCh <- err
				return
			}
		}
	}()
	return conn, nil
}

// Ping sends ping to the TRC and waits for response.
func (c *Conn) Ping() error {
	return c.encoder.Encode(api.NewMessage(api.MessageTypePing, nil, nil))
}

// SetState sends the state to TRC and waits for response.
func (c *Conn) SendState(st *api.State, parent *ulid.ULID) error {
	b, err := json.Marshal(st)
	if err != nil {
		return err
	}
	return c.encoder.Encode(api.NewMessage(api.MessageTypeState, b, parent))
}

// SendHandshake sends handshake message.
func (c *Conn) SendHandshake(hs *api.Handshake) error {
	b, err := json.Marshal(hs)
	if err != nil {
		return err
	}
	return c.encoder.Encode(api.NewMessage(api.MessageTypeHandshake, b, nil))
}

// Close closes the connection.
func (c *Conn) Close() error {
	close(c.closeCh)
	return nil
}

// Errors returns a channel, on which errors are sent.
// There should be exactly one goroutine reading on the returned channel at all times.
func (c *Conn) Errors() <-chan error {
	return c.errCh
}