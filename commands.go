package turtlitto

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
	BallFoundYes           BallFound = "yes"
	BallFoundCommunicatied BallFound = "communicated"
	BallFoundNo            BallFound = "no"
)

type CPB string

const (
	CPBYes           CPB = "yes"
	CPBCommunicatied CPB = "team"
	CPBNo            CPB = "no"
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
