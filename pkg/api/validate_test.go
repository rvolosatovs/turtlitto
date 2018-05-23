package api_test

import (
	"strconv"
	"testing"

	. "github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type TestCase struct {
		Input       Validator
		ShouldError bool
	}

	for i, tc := range []TestCase{
		{ // 0
			Input: &State{
				VisionStatus: "Manual",
			},
			ShouldError: true,
		},
		{ // 1
			Input: &State{
				Temperature1:    28,
				CamStatus:       12, // too much
				Cpu1Load:        89,
				Cpu0Load:        87,
				EmergencyStatus: 0,
			},
			ShouldError: true,
		},
		{ // 2
			Input: &State{
				EmergencyStatus: 100, // maximum
			},
			ShouldError: false,
		},
		{ // 3
			Input: &State{
				BatteryVoltage: 255, // max_uint8
			},
			ShouldError: true,
		},
		{ // 4
			Input: &State{
				ID: "", // empty ID
			},
			ShouldError: false, // not anymore
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := tc.Input.Validate()
			if tc.ShouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
