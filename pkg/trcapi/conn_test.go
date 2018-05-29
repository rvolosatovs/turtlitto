package trcapi_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"context"

	"github.com/rvolosatovs/turtlitto/pkg/api"
	. "github.com/rvolosatovs/turtlitto/pkg/trcapi"
	"github.com/rvolosatovs/turtlitto/pkg/trcapi/trctest"
	"github.com/stretchr/testify/require"
)

func TestState(t *testing.T) {
	for i, tc := range []struct {
		Expected map[string]*api.TurtleState
	}{
		{
			Expected: map[string]*api.TurtleState{
				"foo": {
					BatteryVoltage: 42,
				},
				"bar": {
					HomeGoal: api.HomeGoalBlue,
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

			stateHandler := func(m *api.Message) (*api.Message, error) {
				writes++

				a.Nil(m.Payload)
				pld, err := json.Marshal(tc.Expected)
				if err != nil {
					return nil, err
				}

				return &api.Message{
					MessageID: m.MessageID,
					Type:      m.Type,
					Payload:   pld,
				}, nil
			}

			trctest.Connect(out, in,
				&api.Handshake{Version: DefaultVersion},
				trctest.WithHandler(MessageTypeState, stateHandler),
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
		Input  map[string]*api.TurtleState
		Output map[string]*api.TurtleState
	}{
		{
			Input: map[string]*api.TurtleState{
				"foo": {
					CPB: api.CPBNo,
				},
				"bar": {
					RefBoxRole:   api.RoleDefenderAssist,
					Kinect2State: api.KinectStateBall,
				},
				"baz": {
					ActiveDevPC:    uint8(42),
					BatteryVoltage: uint8(69),
				},
			},
			Output: map[string]*api.TurtleState{
				"foo": {
					CPB: api.CPBYes,
				},
				"bar": {
					RefBoxRole:   api.RoleDefenderAssist2,
					Kinect2State: api.KinectStateNoBall,
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

			state := &api.State{Turtles: tc.Input}
			s := conn.SetState(context.Background(), state)

			a.Nil(err)
			a.Equal(tc.Output, s)
			a.Equal(writes, 1)
		})
	}
}

func TestCommand(t *testing.T) {
	for i, tc := range []struct {
		Command api.Command
	}{
		{
			Command: api.CommandStop,
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
