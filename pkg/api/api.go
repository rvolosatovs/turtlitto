package api

import (
	"crypto/rand"
	"encoding/json"

	"github.com/blang/semver"
	"github.com/oklog/ulid"
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

type RefBoxRole string

const (
	RefBoxRole1 RefBoxRole = "role_1"
	RefBoxRole2 RefBoxRole = "role_2"
	RefBoxRole3 RefBoxRole = "role_3"
	RefBoxRole4 RefBoxRole = "role_4"
	RefBoxRole5 RefBoxRole = "role_5"
	RefBoxRole6 RefBoxRole = "role_6"
)

type KinectState string

const (
	KinectStateNoState KinectState = "no_state"
	KinectStateNoBall  KinectState = "no_ball"
	KinectStateBall    KinectState = "ball"
)

type LocalizationStatus string

const (
	LocalizationStatusNoLocalization LocalizationStatus = "no_localization"
	LocalizationStatusCompassError   LocalizationStatus = "compass_error"
	LocalizationStatusLocalization   LocalizationStatus = "localization"
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

// TurtleState is the state of a particular turtle.
type TurtleState struct {
	// VisionStatus represents status of Vision Executable.
	VisionStatus *bool `json:"visionstatus,omitempty"`

	// MotionStatus represents status of Motion Executable (Off/On).
	MotionStatus *bool `json:"motionstatus,omitempty"`

	// WorldmodelStatus represents status of Worldmodel Executable (Off/On).
	WorldmodelStatus *bool `json:"worldmodelstatus,omitempty"`

	// AppmanStatus represents status of Appman (Off/On).
	AppmanStatus *bool `json:"appmanstatus,omitempty"`

	// RestartCountMotion represents restart count of Motion Executable (0 … 99).
	RestartCountMotion uint8 `json:"restartcountmotion,omitempty"`

	// RestartCountVision represents restart count of Vision Executable (0 … 99).
	RestartCountVision uint8 `json:"restartcountvision,omitempty"`

	// RestartCountWorldmodel represents restart count of Worldmodel Executable (0 … 99).
	RestartCountWorldmodel uint8 `json:"restartcountworldmodel,omitempty"`

	// BallFound represents ball Found (No/Communicated/Yes).
	BallFound BallFound `json:"ballfound,omitempty"`

	// LocalizationStatus represents localization Status.
	LocalizationStatus LocalizationStatus `json:"localizationstatus,omitempty"`

	// CPB represents current Ball Possessor (No/Team/Yes).
	CPB CPB `json:"cpb,omitempty"`

	// BatteryVoltage represents battery Voltage (0 … 99).
	BatteryVoltage uint8 `json:"batteryvoltage,omitempty"`

	// EmergencyStatus represents emergency Status (0 100).
	EmergencyStatus uint8 `json:"emergencystatus,omitempty"`

	// Role represents TRC Role (0 … 10).
	Role Role `json:"role,omitempty"`

	// RefBoxRole represents TRC RefboxRole (0 … 10).
	RefBoxRole RefBoxRole `json:"refboxrole,omitempty"`

	// RobotInField represents TRC Robot In Field (0/1).
	RobotInField *bool `json:"robotinfield,omitempty"`

	// RobotEmergencyButton represents TRC Robot Emergency Button pressed (0/1).
	RobotEmergencyButton *bool `json:"robotembutton,omitempty"`

	// HomeGoal represents robot’s HomeGoal (Yellow/Blue).
	HomeGoal HomeGoal `json:"homegoal,omitempty"`

	// TeamColor represents robot’s Teamcolor (Magenta/Cyan).
	TeamColor TeamColor `json:"teamcolor,omitempty"`

	// ActiveDevPC represents active DevPC controlling robot (0 … 90).
	ActiveDevPC uint8 `json:"activedevpc,omitempty"`

	// Kinect1State represents status of Kinect 1 (No State/No Ball/Ball).
	Kinect1State KinectState `json:"kinect1_state,omitempty"`

	// Kinect2State represents status of Kinect 2 (No State/No Ball/Ball).
	Kinect2State KinectState `json:"kinect2_state,omitempty"`
}

// Message specifies the type of the message.
type MessageType string

const (
	MessageTypeState     MessageType = "state"
	MessageTypePing      MessageType = "ping"
	MessageTypeHandshake MessageType = "handshake"
)

// Handshake represents the handshake message payload.
type Handshake struct {
	Version semver.Version `json:"version"`
	Token   string         `json:"token"`
}

// State represents the state of the TRC.
type State struct {
	Command Command                 `json:"command,omitempty"`
	Turtles map[string]*TurtleState `json:"turtles,omitempty"`
}

// Message is the structure exchanged between TRC and SRRS.
type Message struct {
	Type      MessageType     `json:"type"`
	MessageID ulid.ULID       `json:"message_id"`
	ParentID  *ulid.ULID      `json:"parent_id,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// NewMessage returns a new Message.
func NewMessage(typ MessageType, pld json.RawMessage, parentID *ulid.ULID) *Message {
	return &Message{
		Type:      typ,
		MessageID: ulid.MustNew(ulid.Now(), rand.Reader),
		ParentID:  parentID,
		Payload:   pld,
	}
}
