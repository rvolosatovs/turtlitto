package api

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func Example(t *testing.T) {
	b, err := json.Marshal(map[string]*State{
		"foo": &State{
			Role: "goalkeeper",
		},
		"bar": &State{
			Role: "goalsaver",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
