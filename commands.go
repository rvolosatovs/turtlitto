package turtlitto

// TRCState is the state of the TRC.
type TRCState struct {
	Teams   map[string]*TeamState   `json:"teams"`
	Turtles map[string]*TurtleState `json:"turtles"`
	// TODO: Define types, json tags
}

type TeamState struct {
	Mode string `json:"mode"`
	// TODO: Define types, json tags
}

// TurtleState is the state of a particular turtle.
type TurtleState struct {
	Role string `json:"role"`
	// TODO: Define types, json tags
}

// TurtleStatus is the status of a particular turtle.
type TurtleStatus struct {
	ID      string `json:"id"`
	Battery uint   `json:"battery"`
	// TODO: Define types, json tags
	//motionstatus
	//visionstatus
	//worldmodelstatus
	//appmanstatus
	//commstatus
	//sofsvnrev
	//restartcountmotion
	//restartcountvision
	//restartcountworldmodel
	//ballfound
	//localizationstatus
	//cpb
	//batteryvoltage
	//emergencystatus
	//cpu0load
	//cpu1load
	//role
	//refboxrole
	//robotinfield
	//robotembutton
	//homegoal
	//teamcolor
	//activedevpc
	//temperature_m1
	//temperature_m2
	//temperature_m3
	//is_active
	//cam_status
	//libsvnrev
	//capacitorstate
	//kinect1_state
	//kinect2_state
}
