package api

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"sync"

	"github.com/blang/semver"
	"github.com/mohae/deepcopy"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
)

// DefaultVersion represents the default protocol version.
var DefaultVersion = semver.MustParse("1.0.0")

// ErrClosed represents an error, which occurs when the *Conn is closed.
var ErrClosed = errors.New("Conn is closed")

// NewMessage returns a new Message.
func NewMessage(typ MessageType, pld json.RawMessage, parentID *ulid.ULID) *Message {
	return &Message{
		Type:      typ,
		MessageID: ulid.MustNew(ulid.Now(), rand.Reader),
		ParentID:  parentID,
		Payload:   pld,
	}
}

type encoder interface {
	Encode(v interface{}) (err error)
}

type decoder interface {
	Decode(v interface{}) (err error)
}

// Conn is a connection to TRC.
// Conn is safe for concurrent use by multiple goroutines.
type Conn struct {
	version semver.Version

	decoder decoder
	encoder encoder

	closeChMu *sync.RWMutex
	closeCh   chan struct{}

	errCh chan error

	stateMu *sync.RWMutex
	// state is the current state of TRC.
	state *State

	stateSubsMu *sync.RWMutex
	stateSubs   map[chan<- struct{}]struct{}

	pendingReqsMu *sync.RWMutex
	pendingReqs   map[ulid.ULID]chan *Message
}

// Connect establishes the SRRS-side connection according to TRC API protocol
// specification of version ver on w and r.
func Connect(ver semver.Version, w io.Writer, r io.Reader) (*Conn, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	conn := &Conn{
		version:       ver,
		closeChMu:     &sync.RWMutex{},
		closeCh:       make(chan struct{}),
		decoder:       dec,
		encoder:       json.NewEncoder(w),
		errCh:         make(chan error),
		stateMu:       &sync.RWMutex{},
		state:         &State{},
		stateSubsMu:   &sync.RWMutex{},
		stateSubs:     make(map[chan<- struct{}]struct{}),
		pendingReqsMu: &sync.RWMutex{},
		pendingReqs:   make(map[ulid.ULID]chan *Message),
	}

	var req Message
	if err := conn.decoder.Decode(&req); err != nil {
		return nil, err
	}

	if req.Type != MessageTypeHandshake {
		return nil, errors.Errorf("Expected message of type %s, got %s", MessageTypeHandshake, req.Type)
	}
	if len(req.Payload) == 0 {
		return nil, errors.New("Handshake payload is empty")
	}

	var hs Handshake
	if err := json.Unmarshal(req.Payload, &hs); err != nil {
		return nil, err
	}

	resp := &Handshake{
		Version: hs.Version,
	}
	switch {
	case resp.Version.Major != ver.Major:
		return nil, errors.New("Major version mismatch")
	case resp.Version.Minor > ver.Minor:
		resp.Version = ver
	}
	conn.version = resp.Version

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	if err := conn.encoder.Encode(NewMessage(req.Type, b, &req.MessageID)); err != nil {
		return nil, err
	}

	go func() {
		for {
			var msg Message
			err := conn.decoder.Decode(&msg)

			select {
			case <-conn.closeCh:
				// Don't handle err if connection is closed
				close(conn.errCh)
				return
			}
			if err != nil {
				conn.errCh <- err
				return
			}

			switch msg.Type {
			case MessageTypePing:
				if msg.ParentID != nil {
					// Don't respond to a pong
					break
				}

				if err := conn.encoder.Encode(NewMessage(MessageTypePing, nil, &msg.MessageID)); err != nil {
					conn.errCh <- err
					return
				}

			case MessageTypeState:
				conn.stateMu.Lock()
				st := deepcopy.Copy(conn.state).(*State)
				if err := json.Unmarshal(msg.Payload, st); err != nil {
					conn.stateMu.Unlock()
					conn.errCh <- err
					return
				}
				conn.state = st
				conn.stateMu.Unlock()

				conn.stateSubsMu.RLock()
				for ch := range conn.stateSubs {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
				conn.stateSubsMu.RUnlock()
			default:
				conn.errCh <- errors.Errorf("Unmatched message type: %s", msg.Type)
				return
			}

			conn.pendingReqsMu.RLock()
			ch, ok := conn.pendingReqs[msg.MessageID]
			if ok {
				ch <- &msg
				close(ch)
			}
			conn.pendingReqsMu.RUnlock()
		}
	}()
	return conn, nil
}

// sendRequest sends a request of type typ with payload pld and waits for the response.
func (c *Conn) sendRequest(ctx context.Context, typ MessageType, pld interface{}) (json.RawMessage, error) {
	c.closeChMu.RLock()
	defer c.closeChMu.RUnlock()

	select {
	case <-c.closeCh:
		return nil, ErrClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	v, ok := pld.(Validator)
	if ok && v != nil && v.Validate() != nil {
		return nil, errors.Wrap(v.Validate(), "payload is invalid")
	}

	b, err := json.Marshal(pld)
	if err != nil {
		return nil, err
	}

	msg := NewMessage(typ, b, nil)

	ch := make(chan *Message, 1)
	c.pendingReqsMu.Lock()
	c.pendingReqs[msg.MessageID] = ch
	c.pendingReqsMu.Unlock()

	if err := c.encoder.Encode(msg); err != nil {
		return nil, err
	}

	resp := <-ch

	c.pendingReqsMu.Lock()
	delete(c.pendingReqs, msg.MessageID)
	c.pendingReqsMu.Unlock()
	return resp.Payload, nil
}

// Close closes the connection.
func (c *Conn) Close() error {
	close(c.closeCh)
	c.stateSubsMu.Lock()
	for ch := range c.stateSubs {
		delete(c.stateSubs, ch)
		close(ch)
	}
	c.stateSubsMu.Unlock()
	return nil
}

// State returns the current state of TRC and turtles.
func (c *Conn) State(_ context.Context) *State {
	c.stateMu.RLock()
	st := deepcopy.Copy(c.state).(*State)
	c.stateMu.RUnlock()
	return st
}

// SubscribeStateChanges opens a subscription to state changes.
// SubscribeStateChanges returns read-only channel, on which a value is sent
// every time there is a state change and a function, which must be used to close the subscription.
func (c *Conn) SubscribeStateChanges(ctx context.Context) (<-chan struct{}, func(), error) {
	c.closeChMu.RLock()
	defer c.closeChMu.RUnlock()

	select {
	case <-c.closeCh:
		return nil, nil, ErrClosed
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}

	c.stateSubsMu.Lock()
	ch := make(chan struct{}, 1)
	c.stateSubs[ch] = struct{}{}
	c.stateSubsMu.Unlock()

	return ch, func() {
		c.stateSubsMu.Lock()
		delete(c.stateSubs, ch)
		c.stateSubsMu.Unlock()

		for {
			// Drain channel
			select {
			case <-ch:
			default:
				close(ch)
				return
			}
		}
	}, nil
}

// Ping sends ping to the TRC and waits for response.
func (c *Conn) Ping(ctx context.Context) error {
	_, err := c.sendRequest(ctx, MessageTypePing, nil)
	return err
}

// SetState sends the state to TRC and waits for response.
func (c *Conn) SetState(ctx context.Context, st *State) error {
	_, err := c.sendRequest(ctx, MessageTypeState, st)
	return err
}

// SetCommand sends a command to TRC and waits for response.
func (c *Conn) SetCommand(ctx context.Context, cmd Command) error {
	return c.SetState(ctx, &State{
		Command: cmd,
	})
}

// SetTurtleState sends a state of particular turtle to TRC and waits for response.
func (c *Conn) SetTurtleState(ctx context.Context, id string, st *TurtleState) error {
	return c.SetState(ctx, &State{
		Turtles: map[string]*TurtleState{
			id: st,
		},
	})
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
