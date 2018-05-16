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
		Expected map[string]*State
	}{
		{
			Expected: map[string]*State{
				"foo": {
					ID:             "foo",
					BatteryVoltage: 42,
				},
				"bar": {
					ID:       "foo",
					HomeGoal: HomeGoalBlue,
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assert.New(t)

			out := &bytes.Buffer{}
			in := &mockWriter{
				WriteFunc: func(b []byte) (int, error) {
					var m Message
					err := json.Unmarshal(b, &m)
					a.Nil(err)
					a.NotEmpty(m.MessageID)
					a.Equal(MessageTypeGetState, m.Type)
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

			s, err := NewClient(in, out).State()
			a.Nil(err)
			a.Equal(s, tc.Expected)
		})
	}
}
