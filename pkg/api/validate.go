package api

import (
	"reflect"

	"github.com/pkg/errors"
)

// Validator represents an entity, which can validate itself.
type Validator interface {
	Validate() error
}

// Validate implements Validator.
func (v BallFound) Validate() error {
	switch v {
	case BallFoundYes, BallFoundNo, BallFoundCommunicated:
	default:
		return errors.Errorf("invalid BallFound: %s", v)
	}
	return nil
}

// Validate implements Validator.
func (v LocalizationStatus) Validate() error {
	switch v {
	case LocalizationStatusOn, LocalizationStatusOff, LocalizationStatusManual:
	default:
		return errors.Errorf("invalid LocalizationStatus: %s", v)
	}
	return nil
}

// Validate implements Validator.
func (v CPB) Validate() error {
	switch v {
	case CPBNo, CPBYes, CPBCommunicated:
	default:
		return errors.Errorf("invalid CPB: %s", v)
	}
	return nil
}

// Validate implements Validator.
func (v Role) Validate() error {
	switch v {
	case RoleNone, RoleAttackerMain, RoleAttackerAssist, RoleDefenderMain,
		RoleDefenderAssist, RoleDefenderAssist2, RoleGoalkeeper, RoleInactive:
	default:
		return errors.Errorf("invalid Role: %s", v)
	}
	return nil
}

// Validate implements Validator.
func (v HomeGoal) Validate() error {
	switch v {
	case HomeGoalBlue, HomeGoalYellow:
	default:
		return errors.Errorf("invalid HomeGoal: %s", v)
	}
	return nil
}

// Validate implements Validator.
func (v TeamColor) Validate() error {
	switch v {
	case TeamColorCyan, TeamColorMagenta:
	default:
		return errors.Errorf("invalid TeamColor: %s", v)
	}
	return nil
}

// Validate implements Validator.
func (v KinectState) Validate() error {
	switch v {
	case KinectStateBall, KinectStateNoBall, KinectStateNoState:
	default:
		return errors.Errorf("invalid KinectState: %s", v)
	}
	return nil
}

// rangeError returns an out-of-range error.
func rangeError(source string) error {
	return errors.Errorf("%s out of range", source)
}

// Validate implements Validator.
func (s *TurtleState) Validate() error {
	switch {
	case s.RestartCountVision > 99:
		return rangeError("RestartCountVision")
	case s.RestartCountWorldmodel > 99:
		return rangeError("RestartCountWorldmodel")
	case s.BatteryVoltage > 99:
		return rangeError("BatteryVoltage")
	case s.EmergencyStatus > 100:
		return rangeError("EmergencyStatus")
	case s.ActiveDevPC > 90:
		return rangeError("ActiveDevPC")
	}

	rv := reflect.Indirect(reflect.ValueOf(s))
	for i := 0; i < rv.NumField(); i++ {
		fv := reflect.Indirect(rv.Field(i))
		iface := fv.Interface()
		if !fv.IsValid() || reflect.DeepEqual(iface, reflect.Zero(fv.Type()).Interface()) {
			continue
		}

		v, ok := iface.(Validator)
		if !ok {
			continue
		}

		if err := v.Validate(); err != nil {
			return errors.Wrapf(err, "invalid value of %s", rv.Type().Field(i).Name)
		}
	}
	return nil
}

// Validate implements Validator.
func (s *State) Validate() error {
	if s.Command != "" && s.Validate() != nil { //TODO fix this recursion
		return s.Validate()
	}
	for _, ts := range s.Turtles {
		if err := ts.Validate(); err != nil {
			return err
		}
	}
	return nil
}
