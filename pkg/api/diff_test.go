package api_test

import (
	"strconv"
	"testing"

	. "github.com/rvolosatovs/turtlitto/pkg/api"
	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/stretchr/testify/assert"
)

func TestStateDiff(t *testing.T) {
	type TestCase struct {
		A    *State
		B    *State
		Diff *State
	}

	for i, tc := range []TestCase{
		//test case 1
		{
			A: &State{
				Command: CommandStop,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
					"2": {
						Kinect1State:    KinectStateBall,
						Kinect2State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(1),
					},
					"3": {
						HomeGoal:        HomeGoalYellow,
						BatteryVoltage:  apitest.Uint8Ptr(29),
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			B: &State{
				Command: CommandPenaltyMagenta,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
					"2": {
						Kinect1State:    KinectStateNoBall,
						Kinect2State:    KinectStateNoState,
						EmergencyStatus: apitest.Uint8Ptr(12),
					},
					"3": {
						HomeGoal:        HomeGoalBlue,
						BatteryVoltage:  apitest.Uint8Ptr(23),
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			Diff: &State{
				Command: CommandPenaltyMagenta,
				Turtles: map[string]*TurtleState{
					"2": {
						Kinect1State:    KinectStateNoBall,
						Kinect2State:    KinectStateNoState,
						EmergencyStatus: apitest.Uint8Ptr(12),
					},
					"3": {
						HomeGoal:       HomeGoalBlue,
						BatteryVoltage: apitest.Uint8Ptr(23),
					},
				},
			},
		},
		//test case 2
		{
			A: &State{
				Command: CommandStop,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			B: &State{
				Command: CommandPenaltyMagenta,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			Diff: &State{
				Command: CommandPenaltyMagenta,
			},
		},
		//test case 3
		{
			A: &State{
				Command: CommandStop,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			B: &State{
				Command: CommandStop,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			Diff: nil,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := assert.New(t)

			stateDiff, err := StateDiff(tc.A, tc.B)
			a.Nil(err)
			a.Equal(tc.Diff, stateDiff)
		})
		//t.Run("fuzz", func(t *testing.T) {
		//for i := 0; i < 100; i++ {
		//t.Run(strconv.Itoa(i), func(t *testing.T) {
		//a := assert.New(t)

		//var err error
		//_, err = StateDiff(apitest.RandomState(), apitest.RandomState())
		//a.Nil(err)
		//})
		//}
		//})
	}
}
