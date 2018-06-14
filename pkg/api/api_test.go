package api

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/rvolosatovs/turtlitto/pkg/api/apitest"
	"github.com/stretchr/testify/assert"
)

func TestStateDifference(t *testing.T) {
	type TestCase struct {
		inputOldState       *State
		inputNewState       *State
		expectedOutputState *State
	}

	for i, tc := range []TestCase{
		//test case 1
		{
			inputOldState: &State{
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
			inputNewState: &State{
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
			expectedOutputState: &State{
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
			inputOldState: &State{
				Command: CommandStop,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			inputNewState: &State{
				Command: CommandPenaltyMagenta,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			expectedOutputState: &State{
				Command: CommandPenaltyMagenta,
			},
		},
		//test case 3
		{
			inputOldState: &State{
				Command: CommandStop,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			inputNewState: &State{
				Command: CommandStop,
				Turtles: map[string]*TurtleState{
					"1": {
						Kinect1State:    KinectStateNoBall,
						EmergencyStatus: apitest.Uint8Ptr(0),
					},
				},
			},
			expectedOutputState: &State{},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stateDiff, err := StateDiff(tc.inputOldState, tc.inputNewState)
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(tc.expectedOutputState, stateDiff))
		})
	}
}
