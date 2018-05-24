package api

import (
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
type Conn struct {
	version semver.Version

	decoder decoder
	encoder encoder

	closeCh chan struct{}
	errCh   chan error

	stateMu *sync.RWMutex
	// state is the current state of TRC.
	state *State

	stateSubsMu *sync.RWMutex
	stateSubs   map[chan<- *State]struct{}

	pendingReqsMu *sync.RWMutex
	pendingReqs   map[ulid.ULID]chan struct{}
}

func (c *Conn) sendRequest(typ MessageType, pld interface{}) error {
	v, ok := pld.(Validator)
	if ok && v.Validate() != nil {
		return errors.Wrap(v.Validate(), "payload is invalid")
	}

	b, err := json.Marshal(pld)
	if err != nil {
		return err
	}
	msg := NewMessage(typ, b, nil)

	ch := make(chan struct{})
	c.pendingReqsMu.Lock()
	c.pendingReqs[msg.MessageID] = ch
	c.pendingReqsMu.Unlock()

	if err := c.encoder.Encode(msg); err != nil {
		return err
	}

	<-ch

	c.pendingReqsMu.Lock()
	delete(c.pendingReqs, msg.MessageID)
	c.pendingReqsMu.Unlock()
	return nil
}

func Connect(ver semver.Version, w io.Writer, r io.Reader) (*Conn, error) {
	dec := json.NewDecoder(r)

	var req Message
	if err := dec.Decode(&req); err != nil {
		return nil, err
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

	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(NewMessage(req.Type, b, &req.MessageID)); err != nil {
		return nil, err
	}

	conn := &Conn{
		version:       resp.Version,
		closeCh:       make(chan struct{}),
		errCh:         make(chan error),
		stateMu:       &sync.RWMutex{},
		state:         &State{},
		stateSubsMu:   &sync.RWMutex{},
		stateSubs:     make(map[chan<- *State]struct{}),
		pendingReqsMu: &sync.RWMutex{},
		pendingReqs:   make(map[ulid.ULID]chan struct{}),
	}

	go func() {
		for {
			var msg Message
			err := dec.Decode(&msg)

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
				if err := enc.Encode(NewMessage(MessageTypePing, nil, &msg.MessageID)); err != nil {
					conn.errCh <- err
					return
				}

			case MessageTypeState:
				conn.stateMu.Lock()
				st := deepcopy.Copy(conn.state).(*State)
				if err := json.Unmarshal(msg.Payload, &st); err != nil {
					conn.errCh <- err
					return
				}
				conn.state = st
				conn.stateMu.Unlock()

				conn.stateSubsMu.RLock()
				for ch := range conn.stateSubs {
					ch <- st
				}
				conn.stateSubsMu.RUnlock()
			default:
				conn.errCh <- errors.Errorf("Unmatched message type: %s", msg.Type)
				return
			}

			conn.pendingReqsMu.RLock()
			ch, ok := conn.pendingReqs[msg.MessageID]
			if ok {
				close(ch)
			}
			conn.pendingReqsMu.RUnlock()
		}
	}()
	return conn, nil
}

func (c *Conn) Close() error {
	close(c.closeCh)
	return nil
}

// State returns the current state of turtles.
func (c *Conn) State() *State {
	c.stateMu.RLock()
	st := deepcopy.Copy(c.state).(*State)
	c.stateMu.RUnlock()
	return st
}

func (c *Conn) SubscribeState() (<-chan *State, func()) {
	ch := make(chan *State, 1)
	c.stateSubsMu.Lock()
	c.stateSubs[ch] = struct{}{}
	c.stateSubsMu.Unlock()
	return ch, func() {
		c.stateSubsMu.Lock()
		delete(c.stateSubs, ch)
		c.stateSubsMu.Unlock()
		close(ch)
	}
}

func (c *Conn) SetState(st *State) error {
	return c.sendRequest(MessageTypeState, st)
}

func (c *Conn) SetCommand(cmd Command) error {
	return c.SetState(&State{
		Command: cmd,
	})
}

func (c *Conn) SetTurtleState(id string, st *TurtleState) error {
	return c.SetState(&State{
		Turtles: map[string]*TurtleState{
			id: st,
		},
	})
}

func (c *Conn) Errors() chan<- error {
	return c.errCh
}
