// Package trcapi implements the TRC API as defined in the specification.
package trcapi

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"

	"github.com/blang/semver"
	"github.com/mohae/deepcopy"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/logcontext"
	"go.uber.org/zap"
)

// DefaultVersion represents the default protocol version.
var DefaultVersion = semver.MustParse("1.0.0")

// ErrClosed represents an error, which occurs when the *Conn is closed.
var ErrClosed = errors.New("Conn is closed")

// encoder encodes values.
type encoder interface {
	Encode(v interface{}) (err error)
}

// decoder decodes values.
type decoder interface {
	Decode(v interface{}) (err error)
}

// Conn is a connection to TRC.
// Conn is safe for concurrent use by multiple goroutines.
type Conn struct {
	version semver.Version
	token   *atomic.Value

	decoder decoder
	encoder encoder

	closeChMu *sync.RWMutex
	closeCh   chan struct{}

	errCh chan error

	stateMu *sync.RWMutex
	// state is the current state of TRC.
	state *api.State

	stateSubsMu *sync.RWMutex
	stateSubs   map[chan<- struct{}]struct{}

	pendingReqsMu *sync.RWMutex
	pendingReqs   map[ulid.ULID]chan *api.Message
}

// Connect establishes the SRRS-side connection according to TRC API protocol
// specification of version ver.
// Messages are written to w and read from r.
func Connect(ver semver.Version, w io.Writer, r io.Reader) (*Conn, error) {
	logger := zap.L()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	conn := &Conn{
		version:   ver,
		token:     &atomic.Value{},
		closeChMu: &sync.RWMutex{},
		closeCh:   make(chan struct{}),
		decoder:   dec,
		encoder:   json.NewEncoder(w),
		errCh:     make(chan error),
		stateMu:   &sync.RWMutex{},
		state: &api.State{
			Turtles: map[string]*api.TurtleState{
				"1": {},
				"2": {},
				"3": {},
				"4": {},
				"5": {},
				"6": {},
			},
		},
		stateSubsMu:   &sync.RWMutex{},
		stateSubs:     make(map[chan<- struct{}]struct{}),
		pendingReqsMu: &sync.RWMutex{},
		pendingReqs:   make(map[ulid.ULID]chan *api.Message),
	}

	var req api.Message
	if err := conn.decoder.Decode(&req); err != nil {
		return nil, errors.Wrap(err, "failed to decode handshake request message")
	}

	if req.ParentID != nil {
		return nil, errors.New("initial message is a response, while a request was expected")
	}

	logger.Debug("Initial message decoded successfully",
		zap.Reflect("msg", req),
	)

	if req.Type != api.MessageTypeHandshake {
		return nil, errors.Errorf("expected message of type %s, got %s", api.MessageTypeHandshake, req.Type)
	}
	if len(req.Payload) == 0 {
		return nil, errors.New("received handshake payload is empty")
	}

	logger.Debug("Decoding handshake payload...")
	var hs api.Handshake
	if err := json.Unmarshal(req.Payload, &hs); err != nil {
		return nil, errors.Wrap(err, "failed to decode handshake")
	}
	logger.Debug("Handshake payload decoded successfully",
		zap.Stringer("version", hs.Version),
	)

	resp := &api.Handshake{
		Version: hs.Version,
	}
	switch {
	case resp.Version.Major != ver.Major:
		return nil, errors.New("major version mismatch")
	case resp.Version.Minor > ver.Minor:
		resp.Version = ver
	}
	conn.version = resp.Version

	logger.Debug("Updating token...")
	conn.token.Store(hs.Token)

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode handshake response message")
	}

	logger.Debug("Encoding handshake response...",
		zap.Stringer("version", resp.Version),
	)
	if err := conn.encoder.Encode(api.NewMessage(req.Type, b, &req.MessageID)); err != nil {
		return nil, err
	}

	go func() {
		for {
			var msg api.Message
			err := conn.decoder.Decode(&msg)
			if err == io.EOF {
				logger.Debug("EOF during decoding - closing error channel, return...")
				close(conn.errCh)
				return
			}

			select {
			case <-conn.closeCh:
				logger.Debug("Conn closed - closing error channel, return...")
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
				zap.Reflect("msg", msg),
			)

			switch msg.Type {
			case api.MessageTypePing:
				if msg.ParentID != nil {
					// Don't respond to a pong
					break
				}

				if err := conn.encoder.Encode(api.NewMessage(api.MessageTypePing, nil, &msg.MessageID)); err != nil {
					conn.errCh <- errors.Wrap(err, "failed to encode ping message")
					continue
				}

			case api.MessageTypeState:
				conn.stateMu.Lock()
				st := deepcopy.Copy(conn.state).(*api.State)
				if err := json.Unmarshal(msg.Payload, st); err != nil {
					conn.stateMu.Unlock()
					conn.errCh <- errors.Wrap(err, "failed to decode state message payload")
					continue
				}

				logger.Debug("Received state update", zap.Reflect("state", st))

				conn.state = st
				conn.stateMu.Unlock()

				conn.stateSubsMu.RLock()
				for ch := range conn.stateSubs {
					select {
					case ch <- struct{}{}:
						logger.Debug("Sending state update notification...")
					default:
						logger.Debug("Skipping state update...")
					}
				}
				conn.stateSubsMu.RUnlock()

			default:
				logger.Error("Received message of unmatched type")
				conn.errCh <- errors.Errorf("unmatched message type: %s", msg.Type)
				return
			}

			if msg.ParentID == nil {
				continue
			}

			conn.pendingReqsMu.RLock()
			ch, ok := conn.pendingReqs[*msg.ParentID]
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
func (c *Conn) sendRequest(ctx context.Context, typ api.MessageType, pld interface{}) (json.RawMessage, error) {
	logger := zap.L()

	c.closeChMu.RLock()
	defer c.closeChMu.RUnlock()

	select {
	case <-c.closeCh:
		return nil, ErrClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	v, ok := pld.(api.Validator)
	if ok && v != nil && v.Validate() != nil {
		return nil, errors.Wrap(v.Validate(), "payload is invalid")
	}

	b, err := json.Marshal(pld)
	if err != nil {
		return nil, err
	}

	msg := api.NewMessage(typ, b, nil)

	logger = logger.With(
		zap.Reflect("msg", msg),
	)

	ch := make(chan *api.Message, 1)
	c.pendingReqsMu.Lock()
	c.pendingReqs[msg.MessageID] = ch
	c.pendingReqsMu.Unlock()
	defer func() {
		c.pendingReqsMu.Lock()
		delete(c.pendingReqs, msg.MessageID)
		c.pendingReqsMu.Unlock()
	}()

	logger.Debug("Sending request to TRC...")
	if err := c.encoder.Encode(msg); err != nil {
		logger.Error("Failed to send request to TRC", zap.Error(err))
		return nil, err
	}

	var resp *api.Message
	select {
	case <-ctx.Done():
		logger.Debug("Context done, cancelling", zap.Error(err))
		return nil, ctx.Err()
	case resp = <-ch:
		logger.Debug("Response received",
			zap.Reflect("resp", resp),
		)
	}
	return resp.Payload, nil
}

// Close closes the connection.
func (c *Conn) Close() error {
	select {
	case <-c.closeCh:
	default:
		close(c.closeCh)
	}

	c.stateSubsMu.Lock()
	for ch := range c.stateSubs {
		delete(c.stateSubs, ch)
		close(ch)
	}
	c.stateSubsMu.Unlock()
	return nil
}

// State returns the current state of TRC and turtles.
func (c *Conn) State(_ context.Context) *api.State {
	c.stateMu.RLock()
	st := deepcopy.Copy(c.state).(*api.State)
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
	_, err := c.sendRequest(ctx, api.MessageTypePing, nil)
	return err
}

// SetState sends the state to TRC and waits for response.
func (c *Conn) SetState(ctx context.Context, st *api.State) error {
	logcontext.Logger(ctx).Debug("Sending state...",
		zap.Reflect("state", st),
	)
	_, err := c.sendRequest(ctx, api.MessageTypeState, st)
	return err
}

// SetCommand sends a command to TRC and waits for response.
func (c *Conn) SetCommand(ctx context.Context, cmd api.Command) error {
	return c.SetState(ctx, &api.State{
		Command: cmd,
	})
}

// SetTurtleState sends a state of particular turtle to TRC and waits for response.
func (c *Conn) SetTurtleState(ctx context.Context, st map[string]*api.TurtleState) error {
	if len(st) == 0 {
		return errors.New("Empty state specified")
	}
	return c.SetState(ctx, &api.State{
		Turtles: st,
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
