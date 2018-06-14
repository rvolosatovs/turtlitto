// Package trctest is used to mock the TRC.
package trctest

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"go.uber.org/zap"
)

// DefaultStateHandler is a state handler, which sends the request state back as response.
func DefaultStateHandler(msg *api.Message) (*api.Message, error) {
	if msg.ParentID == nil {
		return api.NewMessage(api.MessageTypeState, msg.Payload, &msg.MessageID), nil
	}
	return nil, errors.New("TRC should not receive state responses")
}

// DefaultPingHandler is a ping handler, which responds to pongs.
func DefaultPingHandler(msg *api.Message) (*api.Message, error) {
	if msg.ParentID == nil {
		return api.NewMessage(api.MessageTypePing, nil, &msg.MessageID), nil
	}
	return nil, nil
}

// DefaultHandshakeHandler is a handshake handler, which compares the version to trcapi.DefaultVersion.
func DefaultHandshakeHandler(msg *api.Message) (*api.Message, error) {
	if msg.ParentID == nil {
		return nil, errors.New("TRC should not receive a handshake request")
	}
	var hs api.Handshake
	if err := json.Unmarshal(msg.Payload, &hs); err != nil {
		return nil, errors.Wrapf(err, "failed to decode handshake payload")
	}
	if hs.Version.Compare(trcapi.DefaultVersion) > 0 {
		return nil, errors.Errorf("unsupported version received: %s", hs.Version)
	}
	return nil, nil
}

// Handler is a function, which handles a message.
type Handler func(*api.Message) (*api.Message, error)

// Conn represents a connection to SRRS.
type Conn struct {
	decoder interface{ Decode(v interface{}) error }
	encoder interface{ Encode(v interface{}) error }

	errCh   chan error
	closeCh chan struct{}

	handlers      *sync.Map
	defaultHander Handler
}

// Option represents a Conn option.
type Option func(*Conn)

// WithHandler allows to specify a custom handler for Conn.
func WithHandler(t api.MessageType, h Handler) Option {
	return func(c *Conn) {
		_, ok := c.handlers.LoadOrStore(t, h)
		if ok {
			panic(errors.Errorf("handler for message type %s is already registered", t))
		}
	}
}

// WithDefaultHandler allows to specify a default handler for Conn.
func WithDefaultHandler(h Handler) Option {
	return func(c *Conn) {
		if c.defaultHander != nil {
			panic(errors.New("default handler is already set"))
		}
		c.defaultHander = h
	}
}

// Connect establishes the TRC-side connection according to TRC API protocol
// specification of version ver on w and r.
func Connect(w io.Writer, r io.Reader, opts ...Option) *Conn {
	logger := zap.L()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	conn := &Conn{
		decoder:  dec,
		encoder:  json.NewEncoder(w),
		closeCh:  make(chan struct{}),
		errCh:    make(chan error),
		handlers: &sync.Map{},
	}
	for _, opt := range opts {
		opt(conn)
	}

	if conn.defaultHander == nil {
		conn.defaultHander = func(msg *api.Message) (*api.Message, error) {
			return nil, errors.Errorf("unmatched handler for type %s", msg.Type)
		}
	}

	go func() {
		for {
			var msg api.Message
			logger.Debug("TRC decoding message...")
			err := conn.decoder.Decode(&msg)
			if err == io.EOF {
				logger.Debug("EOF while decoding")
				close(conn.errCh)
				return
			}

			select {
			case <-conn.closeCh:
				logger.Info("TRC connection closed, returning...")
				// Don't handle err if connection is closed
				close(conn.errCh)
				return
			default:
			}
			if err != nil {
				conn.errCh <- errors.Wrap(err, "failed to decode incoming message")
				return
			}

			logger := logger.With(
				zap.String("type", string(msg.Type)),
				zap.Stringer("message_id", msg.MessageID),
			)
			if msg.ParentID != nil {
				logger = logger.With(zap.Stringer("parent_id", msg.ParentID))
			}

			var h Handler
			v, ok := conn.handlers.Load(msg.Type)
			if !ok {
				h = conn.defaultHander
			} else {
				h = v.(Handler)
			}

			logger.Debug("Executing handler...")
			resp, err := h(&msg)
			if err != nil {
				logger.With(zap.Error(err)).Debug("Executing handler failed")
				conn.errCh <- err
				return
			}

			if resp == nil {
				continue
			}

			logger.Debug("Sending response to SRRS...",
				zap.String("type", string(resp.Type)),
				zap.Stringer("message_id", resp.MessageID),
				zap.Stringer("parent_id", resp.ParentID),
			)
			if err := conn.encoder.Encode(resp); err != nil {
				conn.errCh <- err
				return
			}
		}
	}()
	return conn
}

// Ping sends ping to the TRC and waits for response.
func (c *Conn) Ping() error {
	return c.encoder.Encode(api.NewMessage(api.MessageTypePing, nil, nil))
}

// SetState sends the state to TRC and waits for response.
func (c *Conn) SendState(st *api.State) error {
	b, err := json.Marshal(st)
	if err != nil {
		return err
	}
	return c.encoder.Encode(api.NewMessage(api.MessageTypeState, b, nil))
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
