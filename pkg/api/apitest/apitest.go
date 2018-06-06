package apitest

import (
	crand "crypto/rand"
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/blang/semver"
	"github.com/oklog/ulid"
	"github.com/rvolosatovs/turtlitto/pkg/api"
)

//RandomMessage returns a message with randomly generated fields, within boundaries of the specifications.
func RandomMessage() api.Message {
	var msg = api.Message{}
	msg.Type = *RandomMessageType()
	msg.MessageID = *RandomULID()

	var pld interface{}
	switch msg.Type {
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
		panic("could not marshall payload ")
	}

	msg.Payload = b
	return msg
}

//RandomMessageType randomly returns one of the possible MessageTypes.
func RandomMessageType() *api.MessageType {
	vals := []api.MessageType{api.MessageTypeState, api.MessageTypePing, api.MessageTypeHandshake}
	return &vals[rand.Intn(len(vals))]
}

//RandomHandshake returns a Handshake with randomly generated version string (within v0.0.0 - v9.9.9).
func RandomHandshake() *api.Handshake {
	//TODO: Proper way of creating handshakes
	vers := strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(10)) + "." + strconv.Itoa(rand.Intn(10))
	ver, err := semver.Make(vers)
	if err != nil {
		panic("could not parse version")
	}
	b := make([]byte, 10+rand.Intn(10))
	rand.Read(b)
	return &api.Handshake{
		Version: ver,
		Token:   string(b),
	}
}

//RandomULID returns a random, valid ULID.
func RandomULID() *ulid.ULID {
	newulid := ulid.MustNew(ulid.Now(), crand.Reader)
	return &newulid
}

//RandomState returns a full state existing of a RandomCommand and a RandomTurtleState.
func RandomState() *api.State {
	var pld api.State
	pld.Command = *RandomCommand()
	tstatemap := make(map[string]*api.TurtleState)
	for i := range rand.Perm(6) {
		tstatemap[strconv.Itoa(i+1)] = RandomTurtleState()
	}
	pld.Turtles = tstatemap
	return &pld
}

//RandomCommand returns one of the possible Commands of type Command at random.
func RandomCommand() *api.Command {
	var cmds = []api.Command{
		api.CommandDroppedBall,
		api.CommandStart,
		api.CommandStop,
		api.CommandGoIn,
		api.CommandGoOut,
		api.CommandKickOffCyan,
		api.CommandKickOffMagenta,
		api.CommandFreeKickCyan,
		api.CommandFreeKickMagenta,
		api.CommandGoalKickCyan,
		api.CommandGoalKickMagenta,
		api.CommandThrowInCyan,
		api.CommandThrowInMagenta,
		api.CommandCornerCyan,
		api.CommandCornerMagenta,
		api.CommandPenaltyCyan,
		api.CommandPenaltyMagenta,
		api.CommandRoleAssignerOff,
		api.CommandRoleAssignerOn,
		api.CommandPassDemo,
		api.CommandPenaltyMode,
		api.CommandBallHandlingDemo}

	cmd := cmds[rand.Intn(len(cmds))]
	return &cmd
}

//RandomTurtleState returns a TurtleState with fields randomly filled within specification boundaries.
func RandomTurtleState() *api.TurtleState {
	var ballfound = []api.BallFound{api.BallFoundCommunicated, api.BallFoundYes, api.BallFoundYes}
	var cpb = []api.CPB{api.CPBCommunicated, api.CPBNo, api.CPBYes}
	var homegoals = []api.HomeGoal{api.HomeGoalBlue, api.HomeGoalYellow}
	var teams = []api.TeamColor{api.TeamColorCyan, api.TeamColorMagenta}
	var kinectstates = []api.KinectState{api.KinectStateBall, api.KinectStateNoBall, api.KinectStateNoState}
	var roles = []api.Role{
		api.RoleAttackerAssist,
		api.RoleAttackerMain,
		api.RoleDefenderAssist,
		api.RoleDefenderAssist2,
		api.RoleDefenderMain,
		api.RoleGoalkeeper,
		api.RoleInactive,
		api.RoleNone}
	var refroles = []api.RefBoxRole{
		api.RefBoxRole1,
		api.RefBoxRole2,
		api.RefBoxRole3,
		api.RefBoxRole4,
		api.RefBoxRole5,
		api.RefBoxRole6}

	var locstat = []api.LocalizationStatus{api.LocalizationStatusCompassError,
		api.LocalizationStatusLocalization,
		api.LocalizationStatusNoLocalization}

	var turtstate api.TurtleState = api.TurtleState{
		randomBool(),
		randomBool(),
		randomBool(),
		randomBool(),
		uint8(rand.Intn(100)),
		uint8(rand.Intn(100)),
		uint8(rand.Intn(100)),
		ballfound[rand.Intn(len(ballfound))],
		locstat[rand.Intn(len(locstat))],
		cpb[rand.Intn(len(cpb))],
		uint8(rand.Intn(100)),
		uint8(rand.Intn(101)),
		roles[rand.Intn(len(roles))],
		refroles[rand.Intn(len(refroles))],
		randomBool(),
		randomBool(),
		homegoals[rand.Intn(len(homegoals))],
		teams[rand.Intn(len(teams))],
		uint8(rand.Intn(100)),
		kinectstates[rand.Intn(len(kinectstates))],
		kinectstates[rand.Intn(len(kinectstates))]}

	return &turtstate
}

func randomBool() *bool {
	val := rand.Int()%2 == 0
	return &val
}

func must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}
