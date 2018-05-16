package api

import (
	"crypto/rand"
	"encoding/json"
	"io"

	"github.com/oklog/ulid"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidResultType occurs when the result type is invalid.
	ErrInvalidResultType = errors.New("Result type is invalid")

	// ErrInvalidResultType occurs when the result is invalid.
	ErrInvalidResult = errors.New("Result is invalid")

	// ErrInvalidMessageID occurs when the message ID is invalid.
	ErrInvalidMessageID = errors.New("Message id is invalid")
)

type encoder interface {
	Encode(v interface{}) (err error)
}

type decoder interface {
	Decode(v interface{}) (err error)
}

// Client handles the communication with the TRC.
type Client struct {
	decoder decoder
	encoder encoder
	entropy io.Reader
}

// NewClient returns a new *Client.
func NewClient(w io.Writer, r io.Reader) *Client {
	return &Client{
		decoder: json.NewDecoder(r),
		encoder: json.NewEncoder(w),
		entropy: rand.Reader,
	}
}

// do performs the command specified by typ with optional pld
// and attempts to decode the contents of received Payload field into v, if v != nil.
func (cl *Client) do(typ MessageType, pld interface{}, v interface{}) (err error) {
	id, err := ulid.New(ulid.Now(), cl.entropy)
	if err != nil {
		return errors.Wrap(err, "failed to generate ULID")
	}

	req := &Message{
		Type:      typ,
		MessageID: id,
	}
	if pld != nil {
		req.Payload, err = json.Marshal(pld)
		if err != nil {
			return errors.Wrap(err, "failed to marshal payload")
		}
	}

	if err := cl.encoder.Encode(req); err != nil {
		return errors.Wrap(err, "failed to encode request")
	}

	resp := &Message{}
	if err := cl.decoder.Decode(resp); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if resp.MessageID != req.MessageID {
		return ErrInvalidMessageID
	}
	if resp.Type != req.Type {
		return ErrInvalidResultType
	}

	if v != nil {
		return json.Unmarshal(resp.Payload, v)
	}
	return nil
}

// SendCommand sends executes a global command on the TRC.
func (cl *Client) SendCommand(c Command) error {
	return cl.do(MessageTypeGetState, c, nil)
}

// SetState sets the state of turtles and returns the current turtle state.
func (cl *Client) SetState(s map[string]*State) (map[string]*State, error) {
	ret := make(map[string]*State, len(s))
	return ret, cl.do(MessageTypeSetState, s, &ret)
}

// State returns the current state of turtles.
func (cl *Client) State() (map[string]*State, error) {
	ret := make(map[string]*State)
	return ret, cl.do(MessageTypeGetState, nil, &ret)
}
