package api

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

func diff(a, b map[string]interface{}) map[string]interface{} {
	if reflect.DeepEqual(a, b) {
		return nil
	}

	keys := make(map[string]struct{})
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	ret := make(map[string]interface{}, len(keys))
	for k := range keys {
		av, bv := a[k], b[k]

		if reflect.DeepEqual(av, bv) {
			continue
		}

		if mo, ok := av.(map[string]interface{}); ok {
			if mn, ok := bv.(map[string]interface{}); ok {
				ret[k] = diff(mo, mn)
				continue
			}
		}

		// By some weird reason mapstructure.Decode can't handle pointers, hence, we indirect the value.
		ret[k] = reflect.Indirect(reflect.ValueOf(bv)).Interface()
	}

	if len(ret) == 0 {
		return nil
	}
	return ret
}

//StateDiff returns a *State, which represents the diff of a and b.
func StateDiff(a, b *State) (*State, error) {
	if a == nil {
		return b, nil
	}

	m := diff(structs.Map(a), structs.Map(b))
	if len(m) == 0 {
		return nil, nil
	}

	v := &State{}
	return v, mapstructure.WeakDecode(m, v)
}
