package api_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/rvolosatovs/turtlitto"
	. "github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	for i, tc := range []struct {
		Expected map[string]*turtlitto.State
	}{
		{
			Expected: map[string]*turtlitto.State{
				"foo": {
					ID:             "foo",
					BatteryVoltage: 42,
				},
				"bar": {
					ID:       "foo",
					HomeGoal: turtlitto.HomeGoalBlue,
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(tc.Expected, tc.Expected)

			var out bytes.Buffer
			in := bytes.NewBufferString("") // TODO: encode tc.Expected and set here

			cl := NewClient(&out, in)
			// TODO: Check Status()
			_ = cl
		})
	}
}
