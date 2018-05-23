package api

import (
	"encoding/json"
	"errors"
	"github.com/oklog/ulid"
	"reflect"
)

type Command string

const (
	CommandDroppedBall      Command = "dropped_ball"
	CommandStart            Command = "start"
	CommandStop             Command = "stop"
	CommandGoIn             Command = "go_in"
	CommandGoOut            Command = "go_out"
	CommandKickOffMagenta   Command = "kick_off_magenta"
	CommandKickOffCyan      Command = "kick_off_cyan"
	CommandFreeKickMagenta  Command = "free_kick_magenta"
	CommandFreeKickCyan     Command = "free_kick_cyan"
	CommandGoalKickMagenta  Command = "goal_kick_magenta"
	CommandGoalKickCyan     Command = "goal_kick_cyan"
	CommandThrowInMagenta   Command = "throw_in_magenta"
	CommandThrowInCyan      Command = "throw_in_cyan"
	CommandCornerMagenta    Command = "corner_magenta"
	CommandCornerCyan       Command = "corner_cyan"
	CommandPenaltyMagenta   Command = "penalty_magenta"
	CommandPenaltyCyan      Command = "penalty_cyan"
	CommandRoleAssignerOn   Command = "role_assigner_on"
	CommandRoleAssignerOff  Command = "role_assigner_off"
	CommandPassDemo         Command = "pass_demo"
	CommandPenaltyMode      Command = "penalty_demo"
	CommandBallHandlingDemo Command = "ball_handling_demo"
)

type TeamColor string

const (
	TeamColorMagenta TeamColor = "magenta"
	TeamColorCyan    TeamColor = "cyan"
)

type HomeGoal string

const (
	HomeGoalYellow HomeGoal = "yellow"
	HomeGoalBlue   HomeGoal = "blue"
)

type Role string

const (
	RoleNone            Role = "none"
	RoleInactive        Role = "inactive"
	RoleGoalkeeper      Role = "goalkeeper"
	RoleAttackerMain    Role = "attacker_main"
	RoleAttackerAssist  Role = "attacker_assist"
	RoleDefenderMain    Role = "defender_main"
	RoleDefenderAssist  Role = "defender_assist"
	RoleDefenderAssist2 Role = "defender_assist2"
)

type CapacitorState string

const (
	CapacitorStateNoState CapacitorState = "no_state"
	CapacitorStateEmpty   CapacitorState = "empty"
	CapacitorStateFull    CapacitorState = "full"
)

type KinectState string

const (
	KinectStateNoState KinectState = "no_state"
	KinectStateNoBall  KinectState = "no_ball"
	KinectStateBall    KinectState = "ball"
)

type VisionStatus string

const (
	VisionStatusOff    VisionStatus = "off"
	VisionStatusManual VisionStatus = "manual"
	VisionStatusOn     VisionStatus = "on"
)

type LocalizationStatus string

const (
	LocalizationStatusOff    LocalizationStatus = "off"
	LocalizationStatusManual LocalizationStatus = "compass_issue"
	LocalizationStatusOn     LocalizationStatus = "on"
)

type BallFound string

const (
	BallFoundYes          BallFound = "yes"
	BallFoundCommunicated BallFound = "communicated"
	BallFoundNo           BallFound = "no"
)

type CPB string

const (
	CPBYes          CPB = "yes"
	CPBCommunicated CPB = "team"
	CPBNo           CPB = "no"
)

// State is the state of a particular turtle.
type State struct {
	// ID represents the ID of the turtle.
	ID string `json:"id"`

	// VisionStatus represents status of Vision Executable (Off/Manual/On).
	VisionStatus VisionStatus `json:"visionstatus"`

	// MotionStatus represents status of Motion Executable (Off/On).
	MotionStatus bool `json:"motionstatus"`

	// WorldmodelStatus represents status of Worldmodel Executable (Off/On).
	WorldmodelStatus bool `json:"worldmodelstatus"`

	// AppmanStatus represents status of Appman (Off/On).
	AppmanStatus bool `json:"appmanstatus"`

	// CommStatus represents cOMM Status (On/Off).
	CommStatus bool `json:"commstatus"`

	// SoftwareSVNRevision represents software Revision (0 … 99999).
	SoftwareSVNRevision uint64 `json:"sofsvnrev"`

	// RestartCountMotion represents restart count of Motion Executable (0 … 99).
	RestartCountMotion uint8 `json:"restartcountmotion"`

	// RestartCountVision represents restart count of Vision Executable (0 … 99).
	RestartCountVision uint8 `json:"restartcountvision"`

	// RestartCountWorldmodel represents restart count of Worldmodel Executable (0 … 99).
	RestartCountWorldmodel uint8 `json:"restartcountworldmodel"`

	// BallFound represents ball Found (No/Communicated/Yes).
	BallFound BallFound `json:"ballfound"`

	// LocalizationStatus represents localization Status (No/Compass Issue/Yes).
	LocalizationStatus LocalizationStatus `json:"localizationstatus"`

	// CPB represents current Ball Possessor (No/Team/Yes).
	CPB CPB `json:"cpb"`

	// BatteryVoltage represents battery Voltage (0 … 99).
	BatteryVoltage uint8 `json:"batteryvoltage"`

	// EmergencyStatus represents emergency Status (0 100).
	EmergencyStatus uint8 `json:"emergencystatus"`

	// Cpu0Load represents load CPU0 (0 … 99).
	Cpu0Load uint8 `json:"cpu0load"`

	// Cpu1Load represents load CPU1 (0 … 99).
	Cpu1Load uint8 `json:"cpu1load"`

	// Role represents TRC Role (0 … 10).
	Role Role `json:"role"`

	// RefBoxRole represents TRC RefboxRole (0 … 10).
	RefBoxRole Role `json:"refboxrole"`

	// RobotInField represents TRC Robot In Field (0/1).
	RobotInField bool `json:"robotinfield"`

	// RobotEmergencyButton represents TRC Robot Emergency Button pressed (0/1).
	RobotEmergencyButton bool `json:"robotembutton"`

	// HomeGoal represents robot’s HomeGoal (Yellow/Blue).
	HomeGoal HomeGoal `json:"homegoal"`

	// TeamColor represents robot’s Teamcolor (Magenta/Cyan).
	TeamColor TeamColor `json:"teamcolor"`

	// ActiveDevPC represents active DevPC controlling robot (0 … 90).
	ActiveDevPC uint8 `json:"activedevpc"`

	// Temperature1 represents temperature Motor 1 (0 … 99).
	Temperature1 uint8 `json:"temperature_m1"`

	// Temperature2 represents temperature Motor 2 (0 … 99).
	Temperature2 uint8 `json:"temperature_m2"`

	// Temperature3 represents temperature Motor 3 (0 … 99).
	Temperature3 uint8 `json:"temperature_m3"`

	// IsActive represents TRC is Active.
	IsActive bool `json:"is_active"`

	// CamStatus represents camera Status (0 … 10).
	CamStatus uint8 `json:"cam_status"`

	// LibSVNRevision represents TRC Library revision (0 … 99999).
	LibSVNRevision uint32 `json:"libsvnrev"`

	// CapacitorState represents capacitor state (No State/Empty/Full).
	CapacitorState CapacitorState `json:"capacitorstate"`

	// Kinect1State represents status of Kinect 1 (No State/No Ball/Ball).
	Kinect1State KinectState `json:"kinect1_state"`

	// Kinect2State represents status of Kinect 2 (No State/No Ball/Ball).
	Kinect2State KinectState `json:"kinect2_state"`

	// TODO: Figure out if IsActive corresponds to software being active or "robot"
	// being active in e.g. penalty mode.
}

// Message specifies the type of the message.
type MessageType string

const (
	MessageTypeGetState MessageType = "get_state"
	MessageTypeSetState MessageType = "set_state"
	MessageTypeCommand  MessageType = "command"
)

// Message is the structure exchanged between TRC and SRRS.
type Message struct {
	Type      MessageType     `json:"type"`
	MessageID ulid.ULID       `json:"message_id"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// Validator types can be checked for having a valid state.
type Validator interface {
	// Validate() checks whether this type is in a valid state.
	Validate() error
}

// checks every field of the state, returning an error when any invalid field is encountered
func (state State) Validate() error {
	for _, field := range []Validator{
		state.VisionStatus,
		state.BallFound,
		state.LocalizationStatus,
		state.CPB,
		state.Role,
		state.RefBoxRole,
		state.HomeGoal,
		state.TeamColor,
		state.CapacitorState,
		state.Kinect1State,
		state.Kinect2State,
	} {
		// the String() method is required, and also works for non-string Validators
		// otherwise pointer comparision is used, and this doesn't work
		if reflect.ValueOf(field).String() != reflect.Zero(reflect.TypeOf(field)).String() {
			if err := field.Validate(); err != nil {
				return err
			}
		}
	}

	switch {
	case state.CamStatus > 10:
		return rangeError("Cam status")
	case state.RestartCountVision > 99:
		return rangeError("Restart counter of Vision Executable")
	case state.RestartCountWorldmodel > 99:
		return rangeError("Restart counter of World Model")
	case state.BatteryVoltage > 99:
		return rangeError("Battery voltage")
	case state.EmergencyStatus > 100:
		return rangeError("Emergency status")
	case state.Cpu0Load > 99:
		return rangeError("Cpu 0 load")
	case state.Cpu1Load > 99:
		return rangeError("Cpu 1 load")
	case state.ActiveDevPC > 90:
		return rangeError("Active dev pc")
	case state.Temperature1 > 99:
		return rangeError("Temperature 1")
	case state.Temperature2 > 99:
		return rangeError("Temperature 2")
	case state.Temperature3 > 99:
		return rangeError("Temperature 3")
	}
	return nil
}

// returns an error, signalling that the given value is out of range.
func rangeError(source string) error {
	return errors.New(source + " out of range")
}

func (vs VisionStatus) Validate() error {
	switch vs {
	case VisionStatusManual, VisionStatusOn, VisionStatusOff:
	default:
		return errors.New("Invalid VisionStatus: got '" + string(vs) + "'")
	}
	return nil
}

func (bf BallFound) Validate() error {
	switch bf {
	case BallFoundYes, BallFoundNo, BallFoundCommunicated:
	default:
		return errors.New("Invalid BallFound: got '" + string(bf) + "'")
	}
	return nil
}

func (ls LocalizationStatus) Validate() error {
	switch ls {
	case LocalizationStatusOn, LocalizationStatusOff, LocalizationStatusManual:
	default:
		return errors.New("Invalid Localisation status: got '" + string(ls) + "'")
	}
	return nil
}

func (cpb CPB) Validate() error {
	switch cpb {
	case CPBNo, CPBYes, CPBCommunicated:
	default:
		return errors.New("Invalid CPB: got '" + string(cpb) + "'")
	}
	return nil
}

func (r Role) Validate() error {
	switch r {
	case RoleNone, RoleAttackerMain, RoleAttackerAssist, RoleDefenderMain,
		RoleDefenderAssist, RoleDefenderAssist2, RoleGoalkeeper, RoleInactive:
	default:
		return errors.New("Invalid Role: got '" + string(r) + "'")
	}
	return nil
}

func (hg HomeGoal) Validate() error {
	switch hg {
	case HomeGoalBlue, HomeGoalYellow:
	default:
		return errors.New("Invalid HomeGoal: got '" + string(hg) + "'")
	}
	return nil
}

func (tc TeamColor) Validate() error {
	switch tc {
	case TeamColorCyan, TeamColorMagenta:
	default:
		return errors.New("Invalid Team color: got '" + string(tc) + "'")
	}
	return nil
}

func (cs CapacitorState) Validate() error {
	switch cs {
	case CapacitorStateEmpty, CapacitorStateFull, CapacitorStateNoState:
	default:
		return errors.New("Invalid Capacitor state: got '" + string(cs) + "'")
	}
	return nil
}

func (ks KinectState) Validate() error {
	switch ks {
	case KinectStateBall, KinectStateNoBall, KinectStateNoState:
	default:
		return errors.New("Invalid Kinect state: got '" + string(ks) + "'")
	}
	return nil
}
