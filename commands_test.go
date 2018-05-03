package turtlitto

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func ExampleTRCState(t *testing.T) {
	b, err := json.Marshal(&TRCState{
		Teams: map[string]*TeamState{
			"foo": &TeamState{
				Mode: "kick-off",
			},
			"bar": &TeamState{
				Mode: "kick-on",
			},
		},
		Turtles: map[string]*TurtleState{
			"foo": &TurtleState{
				Role: "goalkeeper",
			},
			"bar": &TurtleState{
				Role: "goalsaver",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
