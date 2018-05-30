package api_test

import (
	"testing"

	. "github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type TestCase struct {
		Name        string
		Input       Validator
		ShouldError bool
	}

	for _, tc := range []TestCase{
		{
			Name: "a simple correct turtleState",
			Input: &TurtleState{
				EmergencyStatus: 100, // maximum
			},
			ShouldError: false,
		},
		{
			Name: "a simple erroneous turtleState",
			Input: &TurtleState{
				BatteryVoltage: 255, // max_uint8
			},
			ShouldError: true,
		},
		{
			Name: "a one-item State",
			Input: &State{
				Turtles: map[string]*TurtleState{
					"t": {
						BallFound: BallFoundCommunicated,
					},
				},
			},
			ShouldError: false,
		},
		{
			Name: "a multi-item, multi-turtle State",
			Input: &State{
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: 0,
					},
					"2": {
						Kinect1State:    KinectStateBall,
						Kinect2State:    KinectStateNoBall,
						EmergencyStatus: 1,
					},
					"3": {
						HomeGoal:        HomeGoalYellow,
						BatteryVoltage:  28,
						EmergencyStatus: 0,
					},
				},
			},
			ShouldError: false,
		},
		{
			Name: "a one-item wrong EmergencyStatus TurtleState",
			Input: &State{
				Turtles: map[string]*TurtleState{
					"t": {
						EmergencyStatus: 255, // too much
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "a one-item wrong RefBoxRole TurtleState",
			Input: &State{
				Turtles: map[string]*TurtleState{
					"t": {
						RefBoxRole: "wrongRole",
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "a multi-item, multi-turtle wrong State",
			Input: &State{
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: 0,
					},
					"2": {
						Kinect1State:    KinectStateBall,
						Kinect2State:    KinectStateNoBall,
						EmergencyStatus: 1,
					},
					"3": {
						HomeGoal:        HomeGoalYellow,
						BatteryVoltage:  100, // too much
						EmergencyStatus: 0,
					},
				},
			},
			ShouldError: true,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.Input.Validate()
			if tc.ShouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
