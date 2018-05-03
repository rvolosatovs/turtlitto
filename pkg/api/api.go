package api

import (
	"encoding/json"
	"io"
	"math/rand"
	"time"

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
		entropy: rand.New(newLockedSource(rand.NewSource(time.Now().UnixNano()))),
	}
}

// do performs the command specified by typ with optional pld
// and attempts to decode the contents of received Payload field into v, if v != nil.
func (cl *Client) do(typ string, pld []byte, v interface{}) (err error) {
	id, err := ulid.New(ulid.Now(), cl.entropy)
	if err != nil {
		return errors.Wrap(err, "failed to generate ULID")
	}

	type message struct {
		Type      string    `json:"type"`
		MessageID ulid.ULID `json:"message_id"`
		Payload   []byte    `json:"payload,omitempty"`
	}

	req := &message{
		Type:      typ,
		MessageID: id,
		Payload:   pld,
	}
	if err := cl.encoder.Encode(req); err != nil {
		return errors.Wrap(err, "failed to encode command")
	}

	resp := &message{}
	if err := cl.decoder.Decode(resp); err != nil {
		return errors.Wrap(err, "failed to encode result")
	}

	if resp.MessageID != req.MessageID {
		return ErrInvalidMessageID
	}
	if resp.Type != req.Type {
		return ErrInvalidResultType
	}

	if v != nil {
		return cl.decoder.Decode(v)
	}
	return nil
}

// SetState sets the state of TRC and returns the current TRC state.
func (cl *Client) SetState(s *turtlitto.TRCState) (*turtlitto.TRCState, error) {
	ret := &turtlitto.TRCState{}
	return ret, cl.do("set_state", s, ret)
}

// Status returns the current status of turtles as reported by TRC.
func (cl *Client) Status() ([]*turtlitto.TurtleStatus, error) {
	ret := &turtlitto.TurtleStatus{}
	return ret, cl.do("get_status", ret)
}
