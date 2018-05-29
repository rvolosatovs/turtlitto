package test

import (
	"encoding/json"
	"io"

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

type Handler func(*api.Message) (*api.Message, error)

type Conn struct {
	sendCh    chan *api.Message
	errCh     chan error
	handshake *api.Handshake
	handlers  map[api.MessageType]Handler
}

type Option func(*Conn)

func WithHandler(msg api.MessageType, handler Handler) Option {
	return func(c *Conn) {
		c.handlers[msg] = handler
	}
}

func WithHandshake(hs *api.Handshake) Option {
	return func(c *Conn) {
		c.handshake = hs
	}
}

// Connect establishes the TRC-side connection according to TRC API protocol
// specification of version ver on w and r.
func Connect(w io.Writer, r io.Reader, opts ...Option) (*Conn, error) {
	conn := &Conn{
		errCh:  make(chan error, 1),
		sendCh: make(chan *api.Message, 1),
	}
	for _, opt := range opts {
		opt(trc)
	}

	// Send handshake
	b, err := json.Marshal(handshakeMsg)
	if err != nil {
		return nil, err
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(api.NewMessage(
		api.MessageTypeHandshake,
		b,
		nil,
	)); err != nil {
		return nil, err
	}

	var resp api.Message
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Type != api.MessageTypeHandshake {
		return nil, errors.New("Reply was not an handshake")
	}

	var hs api.Handshake
	if err := json.Unmarshal(resp.Payload, &hs); err != nil {
		return nil, err
	}

	go handleIncoming(dec, enc, trc.handlers, trc.errCh)

	// send possible external messages to client
	go func() {
		for msg := range trc.sendCh {
			if err := enc.Encode(msg); err != nil {
				trc.errCh <- err
			}
		}
	}()

	return trc, nil
}

func (c *Conn) Errors() {

}

// routine for handling incoming messages using the specified handlers.
func handleIncoming(dec *json.Decoder, enc *json.Encoder, handlers map[api.MessageType]Handler, errChan chan<- error) {
	for {
		var msg api.Message
		if err := dec.Decode(&msg); err != nil {
			errChan <- errors.Wrap(err, "Could not decode incoming message")
			return
		}

		han, ok := handlers[msg.Type]
		if !ok {
			errChan <- errors.Errorf("Unknown message: %s", msg)
		}

		reply, err := han(&msg)
		if err != nil {
			errChan <- errors.Wrap(err, "Handler error")
		}

		enc.Encode(reply)
	}
}

// Errors returns a channel, on which errors are sent.
// There should be exactly one goroutine reading on the returned channel at all times.
func (c *Conn) Errors() <-chan error {
	return c.errCh
}

// Closed returns a channel that's closed when Conn is closed.
// Successive calls to Closed return the same value.
// There may be multiple goroutines reading on the returned channel.
func (c *Conn) Closed() <-chan struct{} {
	return c.closeCh
}
