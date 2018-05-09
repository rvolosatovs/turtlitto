package api

import (
	"crypto/rand"
	"encoding/json"
	"io"

	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	"github.com/rvolosatovs/turtlitto"
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

func NewClient(w io.Writer, r io.Reader) *Client {
	return &Client{
		decoder: json.NewDecoder(r),
		encoder: json.NewEncoder(w),
		entropy: rand.Reader,
	}
}

type message struct {
	Type      string    `json:"type"`
	MessageID ulid.ULID `json:"message_id"`
	Payload   []byte    `json:"payload,omitempty"`
}

// do performs the command specified by typ with optional pld
// and attempts to decode the contents of received Payload field into v, if v != nil.
func (cl *Client) do(typ string, pld interface{}, v interface{}) (err error) {
	id, err := ulid.New(ulid.Now(), cl.entropy)
	if err != nil {
		return errors.Wrap(err, "failed to generate ULID")
	}

	req := &message{
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

	resp := &message{}
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
func (cl *Client) SendCommand(c turtlitto.Command) error {
	return cl.do("command", c, nil)
}

// SetState sets the state of turtles and returns the current turtle state.
func (cl *Client) SetState(s map[string]*turtlitto.State) (map[string]*turtlitto.State, error) {
	ret := make(map[string]*turtlitto.State, len(s))
	return ret, cl.do("set_state", s, &ret)
}

// Status returns the current status of turtles.
func (cl *Client) Status() (map[string]*turtlitto.State, error) {
	ret := make(map[string]*turtlitto.State)
	return ret, cl.do("get_status", nil, &ret)
}
