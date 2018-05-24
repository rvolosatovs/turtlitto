package api_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	. "github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	WriteFunc func(b []byte) (int, error)
}

func (r *mockWriter) Write(b []byte) (int, error) {
	return r.WriteFunc(b)
}

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
			a := assert.New(t)

			writes := 0

			out := &bytes.Buffer{}
			in := &mockWriter{
				WriteFunc: func(b []byte) (int, error) {
					writes++

					var m Message
					err := json.Unmarshal(b, &m)
					a.Nil(err)
					a.NotEmpty(m.MessageID)
					a.Equal(MessageTypeState, m.Type)
					a.Nil(m.Payload)

					pld, err := json.Marshal(tc.Expected)
					if err != nil {
						panic(err)
					}

					err = json.NewEncoder(out).Encode(&Message{
						MessageID: m.MessageID,
						Type:      m.Type,
						Payload:   pld,
					})
					if err != nil {
						panic(err)
					}
					return len(b), nil
				},
			}

			conn, err := Connect(DefaultVersion, in, out)
			a.Nil(err)
			a.Equal(s, tc.Expected)
			a.Equal(writes, 1)
		})
	}
}

func TestSetState(t *testing.T) {
	for i, tc := range []struct {
		Input  map[string]*State
		Output map[string]*State
	}{
		{
			Input: map[string]*State{
				"foo": {
					ID:  "1",
					CPB: CPBNo,
				},
				"bar": {
					ID:           "2",
					RefBoxRole:   RoleDefenderAssist,
					Kinect2State: KinectStateBall,
				},
				"baz": {
					ID:           "3",
					ActiveDevPC:  uint8(42),
					Temperature1: uint8(69),
				},
			},
			Output: map[string]*State{
				"foo": {
					ID:  "1",
					CPB: CPBYes,
				},
				"bar": {
					ID:           "2",
					RefBoxRole:   RoleDefenderAssist2,
					Kinect2State: KinectStateNoBall,
				},
				"baz": {
					ID:           "3",
					ActiveDevPC:  uint8(89),
					Temperature1: uint8(69),
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assert.New(t)

			writes := 0

			out := &bytes.Buffer{}
			in := &mockWriter{
				WriteFunc: func(b []byte) (int, error) {
					writes++

					var m Message
					err := json.Unmarshal(b, &m)
					a.Nil(err)
					a.NotEmpty(m.MessageID)
					a.Equal(MessageTypeSetState, m.Type)

					var ts map[string]*State
					err = json.Unmarshal(m.Payload, &ts)
					a.Nil(err)
					a.Equal(ts, tc.Input)

					pld, err := json.Marshal(tc.Output)
					if err != nil {
						panic(err)
					}

					err = json.NewEncoder(out).Encode(&Message{
						MessageID: m.MessageID,
						Type:      m.Type,
						Payload:   pld,
					})
					if err != nil {
						panic(err)
					}
					return len(b), nil
				},
			}

			cl := NewClient(in, out)

			s, err := cl.SetState(tc.Input)
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
			a := assert.New(t)

			writes := 0

			out := &bytes.Buffer{}
			in := &mockWriter{
				WriteFunc: func(b []byte) (int, error) {
					writes++

					var m Message
					err := json.Unmarshal(b, &m)
					a.Nil(err)
					a.NotEmpty(m.MessageID)
					a.Equal(m.Type, MessageTypeCommand)

					var payload Command
					err = json.Unmarshal(m.Payload, &payload)
					a.Nil(err)
					a.Equal(tc.Command, payload)

					err = json.NewEncoder(out).Encode(&Message{
						MessageID: m.MessageID,
						Type:      m.Type,
					})
					if err != nil {
						panic(err)
					}
					return len(b), nil
				},
			}

			err := NewClient(in, out).SendCommand(tc.Command)
			a.Nil(err)
			a.Equal(writes, 1)
			a.Equal(s, tc.Expected)

			// TODO: Check things
		})
	}
}
