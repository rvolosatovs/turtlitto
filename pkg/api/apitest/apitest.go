// Package apitest provides utilities for API testing.
package apitest

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/blang/semver"
	"github.com/oklog/ulid"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

// MaxTurtles is the maximum amount of turtles.
const MaxTurtles = 1

// BoolPtr returns v as *bool.
func BoolPtr(v bool) *bool {
	return &v
}

// BoolPtr returns v as *uint8.
func Uint8Ptr(v uint8) *uint8 {
	return &v
}

// RandomCommand returns a random valid api.Command.
func RandomCommand() api.Command {
	switch rand.Intn(22) {
	case 0:
		return api.CommandDroppedBall
	case 1:
		return api.CommandStart
	case 2:
		return api.CommandStop
	case 3:
		return api.CommandGoIn
	case 4:
		return api.CommandGoOut
	case 5:
		return api.CommandKickOffCyan
	case 6:
		return api.CommandKickOffMagenta
	case 7:
		return api.CommandFreeKickCyan
	case 8:
		return api.CommandFreeKickMagenta
	case 9:
		return api.CommandGoalKickCyan
	case 10:
		return api.CommandGoalKickMagenta
	case 11:
		return api.CommandThrowInCyan
	case 12:
		return api.CommandThrowInMagenta
	case 13:
		return api.CommandCornerCyan
	case 14:
		return api.CommandCornerMagenta
	case 15:
		return api.CommandPenaltyCyan
	case 16:
		return api.CommandPenaltyMagenta
	case 17:
		return api.CommandRoleAssignerOff
	case 18:
		return api.CommandRoleAssignerOn
	case 19:
		return api.CommandPassDemo
	case 20:
		return api.CommandPenaltyMode
	case 21:
		return api.CommandBallHandlingDemo
	default:
		panic("unmatched")
	}
}

// RandomBallFound returns a random valid api.BallFound.
func RandomBallFound() api.BallFound {
	switch rand.Intn(3) {
	case 0:
		return api.BallFoundYes
	case 1:
		return api.BallFoundCommunicated
	case 2:
		return api.BallFoundNo
	default:
		panic("unmatched")
	}
}

// RandomLocalizationStatus returns a random valid api.LocalizationStatus.
func RandomLocalizationStatus() api.LocalizationStatus {
	switch rand.Intn(3) {
	case 0:
		return api.LocalizationStatusNoLocalization
	case 1:
		return api.LocalizationStatusCompassError
	case 2:
		return api.LocalizationStatusLocalization
	default:
		panic("unmatched")
	}
}

// RandomCPB returns a random valid api.CPB.
func RandomCPB() api.CPB {
	switch rand.Intn(3) {
	case 0:
		return api.CPBYes
	case 1:
		return api.CPBCommunicated
	case 2:
		return api.CPBNo
	default:
		panic("unmatched")
	}
}

// RandomRole returns a random valid api.Role.
func RandomRole() api.Role {
	switch rand.Intn(8) {
	case 0:
		return api.RoleNone
	case 1:
		return api.RoleInactive
	case 2:
		return api.RoleGoalkeeper
	case 3:
		return api.RoleAttackerMain
	case 4:
		return api.RoleAttackerAssist
	case 5:
		return api.RoleDefenderMain
	case 6:
		return api.RoleDefenderAssist
	case 7:
		return api.RoleDefenderAssist2
	default:
		panic("unmatched")
	}
}

// RandomRefboxRole returns a random valid api.RefboxRole.
func RandomRefBoxRole() api.RefBoxRole {
	switch rand.Intn(6) {
	case 0:
		return api.RefBoxRole1
	case 1:
		return api.RefBoxRole2
	case 2:
		return api.RefBoxRole3
	case 3:
		return api.RefBoxRole4
	case 4:
		return api.RefBoxRole5
	case 5:
		return api.RefBoxRole6
	default:
		panic("unmatched")
	}
}

// RandomHomeGoal returns a random valid api.HomeGoal.
func RandomHomeGoal() api.HomeGoal {
	switch rand.Intn(2) {
	case 0:
		return api.HomeGoalYellow
	case 1:
		return api.HomeGoalBlue
	default:
		panic("unmatched")
	}
}

// RandomTeamColor returns a random valid api.TeamColor.
func RandomTeamColor() api.TeamColor {
	switch rand.Intn(2) {
	case 0:
		return api.TeamColorMagenta
	case 1:
		return api.TeamColorCyan
	default:
		panic("unmatched")
	}
}

// RandomKinectState returns a random valid api.KinectState.
func RandomKinectState() api.KinectState {
	switch rand.Intn(3) {
	case 0:
		return api.KinectStateBall
	case 1:
		return api.KinectStateNoBall
	case 2:
		return api.KinectStateNoState
	default:
		panic("unmatched")
	}
}

// RandomTurtleState returns a random valid *api.TurtleState.
func RandomTurtleState() *api.TurtleState {
	return &api.TurtleState{
		VisionStatus:           BoolPtr(rand.Intn(2) == 0),
		MotionStatus:           BoolPtr(rand.Intn(2) == 0),
		WorldmodelStatus:       BoolPtr(rand.Intn(2) == 0),
		AppmanStatus:           BoolPtr(rand.Intn(2) == 0),
		RestartCountMotion:     Uint8Ptr(uint8(rand.Intn(100))),
		RestartCountVision:     Uint8Ptr(uint8(rand.Intn(100))),
		RestartCountWorldmodel: Uint8Ptr(uint8(rand.Intn(100))),
		BallFound:              RandomBallFound(),
		LocalizationStatus:     RandomLocalizationStatus(),
		CPB:                    RandomCPB(),
		BatteryVoltage:         Uint8Ptr(uint8(rand.Intn(100))),
		EmergencyStatus:        Uint8Ptr(uint8(rand.Intn(101))),
		Role:                   RandomRole(),
		RefBoxRole:             RandomRefBoxRole(),
		RobotInField:           BoolPtr(rand.Intn(2) == 0),
		RobotEmergencyButton:   BoolPtr(rand.Intn(2) == 0),
		HomeGoal:               RandomHomeGoal(),
		TeamColor:              RandomTeamColor(),
		ActiveDevPC:            Uint8Ptr(uint8(rand.Intn(91))),
		Kinect1State:           RandomKinectState(),
		Kinect2State:           RandomKinectState(),
	}
}

// RandomTurtleStateMap returns a random valid *api.TurtleState map.
func RandomTurtleStateMap() map[string]*api.TurtleState {
	ret := map[string]*api.TurtleState{}
	perm := rand.Perm(MaxTurtles)[:rand.Intn(MaxTurtles+1)]
	for _, i := range perm {
		ret[strconv.Itoa(i+1)] = RandomTurtleState()
	}
	return ret
}

// RandomState returns a random valid *api.State.
func RandomState() *api.State {
	var pld api.State
	if rand.Intn(2) == 0 {
		pld.Command = RandomCommand()
	}
	if rand.Intn(2) == 0 {
		pld.Turtles = RandomTurtleStateMap()
	}
	return &pld
}

// RandomVersion returns a random valid *semver.Version.
func RandomVersion() *semver.Version {
	ver := semver.MustParse(
		fmt.Sprintf("%d.%d.%d", rand.Intn(10), rand.Intn(10), rand.Intn(10)),
	)
	return &ver
}

// RandomMessageType returns a random valid api.MessageType.
func RandomMessageType() api.MessageType {
	switch rand.Intn(2) {
	case 0:
		return api.MessageTypeState
	case 1:
		return api.MessageTypePing
	case 2:
		return api.MessageTypeHandshake
	default:
		panic("unmatched")
	}
}

// RandomHandshake returns a random valid *api.Handshake.
func RandomHandshake() *api.Handshake {
	b := make([]byte, 10+rand.Intn(10))
	rand.Read(b)
	return &api.Handshake{
		Version: *RandomVersion(),
		Token:   string(b),
	}
}

// RandomULID returns a random valid *ulid.ULID.
func RandomULID() *ulid.ULID {
	ret := ulid.MustNew(ulid.Now(), crand.Reader)
	return &ret
}

// RandomMessage returns a random valid *api.Message.
func RandomMessage() *api.Message {
	var pld interface{}

	mt := RandomMessageType()
	switch mt {
	case api.MessageTypeHandshake:
		pld = *RandomHandshake()
	case api.MessageTypeState:
		pld = *RandomState()
	case api.MessageTypePing:
		pld = nil
	default:
		panic("unmatched Message type")
	}

	b, err := json.Marshal(pld)
	if err != nil {
		panic("failed to marshall payload ")
	}

	var parentID *ulid.ULID
	if rand.Intn(2) == 0 {
		parentID = RandomULID()
	}

	return &api.Message{
		Type:      mt,
		MessageID: *RandomULID(),
		ParentID:  parentID,
		Payload:   b,
	}
}
