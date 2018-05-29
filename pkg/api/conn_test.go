package api_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"context"
	. "github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/stretchr/testify/require"
)

func TestState(t *testing.T) {
	for i, tc := range []struct {
		Expected map[string]*TurtleState
	}{
		{
			Expected: map[string]*TurtleState{
				"foo": {
					BatteryVoltage: 42,
				},
				"bar": {
					HomeGoal: HomeGoalBlue,
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := require.New(t)

			writes := 0

			// simulates the TRC
			out := &bytes.Buffer{}
			in := &bytes.Buffer{}

			stateHandler := func(m *Message) (*Message, error) {
				writes++

				a.Nil(m.Payload)
				pld, err := json.Marshal(tc.Expected)
				if err != nil {
					return nil, err
				}

				return &Message{
					MessageID: m.MessageID,
					Type:      m.Type,
					Payload:   pld,
				}, nil
			}

			test.Connect(out, in,
				&Handshake{Version: DefaultVersion},
				test.WithHandler(MessageTypeState, stateHandler),
			)

			conn, err := Connect(DefaultVersion, in, out)

			a.Nil(err)
			a.Equal(conn.State(context.Background()), tc.Expected)
			a.Equal(writes, 1)
		})
	}
}

func TestSetState(t *testing.T) {
	for i, tc := range []struct {
		Input  map[string]*TurtleState
		Output map[string]*TurtleState
	}{
		{
			Input: map[string]*TurtleState{
				"foo": {
					CPB: CPBNo,
				},
				"bar": {
					RefBoxRole:   RoleDefenderAssist,
					Kinect2State: KinectStateBall,
				},
				"baz": {
					ActiveDevPC:    uint8(42),
					BatteryVoltage: uint8(69),
				},
			},
			Output: map[string]*TurtleState{
				"foo": {
					CPB: CPBYes,
				},
				"bar": {
					RefBoxRole:   RoleDefenderAssist2,
					Kinect2State: KinectStateNoBall,
				},
				"baz": {
					ActiveDevPC:    uint8(89),
					BatteryVoltage: uint8(69),
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := require.New(t)
			writes := 0

			// simulates the TRC
			toSRRS := &bytes.Buffer{}
			toTRC := &bytes.Buffer{}

			stateHandler := func(m *Message) (*Message, error) {
				var ts map[string]*State
				err := json.Unmarshal(m.Payload, &ts)
				a.Nil(err)
				a.Equal(ts, tc.Input)

				pld, err := json.Marshal(tc.Output)
				if err != nil {
					return nil, err
				}
				writes++

				return &Message{
					MessageID: m.MessageID,
					Type:      m.Type,
					Payload:   pld,
				}, nil
			}

			test.Connect(toTRC, toSRRS,
				&Handshake{Version: DefaultVersion},
				test.WithHandler(MessageTypeState, stateHandler),
			)

			conn, err := Connect(DefaultVersion, toSRRS, toTRC)
			a.Nil(err)

			state := &State{Turtles: tc.Input}
			s := conn.SetState(context.Background(), state)

			a.Nil(err)
			a.Equal(tc.Output, s)
			a.Equal(writes, 1)
		})
	}
}

func TestCommand(t *testing.T) {
	for i, tc := range []struct {
		Command Command
	}{
		{
			Command: CommandStop,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := require.New(t)

			writes := 0

			// simulates the TRC
			out := &bytes.Buffer{}
			in := &bytes.Buffer{}

			stateHandler := func(m *Message) (*Message, error) {
				var payload Command
				err := json.Unmarshal(m.Payload, &payload)
				a.Nil(err)
				a.Equal(tc.Command, payload)

				pld, err := json.Marshal(tc.Command)
				if err != nil {
					return nil, err
				}
				writes++

				return &Message{
					MessageID: m.MessageID,
					Type:      m.Type,
					Payload:   pld,
				}, nil
			}

			test.Connect(out, in,
				&Handshake{Version: DefaultVersion},
				test.WithHandler(MessageTypeState, stateHandler),
			)

			conn, err := Connect(DefaultVersion, in, out)
			a.Nil(err)

			s := conn.SetCommand(context.Background(), tc.Command)

			a.Equal(writes, 1)
			a.Equal(s, tc.Command)
		})
	}
}
