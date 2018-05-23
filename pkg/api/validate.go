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
func (s *State) Validate() error {
	if s == nil {
		return nil
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
			return errors.Wrapf(err, "invalid value of %s: %s", rv.Type().Field(i).Name)
		}
	}

	switch {
	case s.CamStatus > 10:
		return rangeError("CamStatus")
	case s.RestartCountVision > 99:
		return rangeError("RestartCountVision")
	case s.RestartCountWorldmodel > 99:
		return rangeError("RestartCountWorldmodel")
	case s.BatteryVoltage > 99:
		return rangeError("BatteryVoltage")
	case s.EmergencyStatus > 100:
		return rangeError("EmergencyStatus")
	case s.Cpu0Load > 99:
		return rangeError("Cpu0Load")
	case s.Cpu1Load > 99:
		return rangeError("Cpu1Load")
	case s.ActiveDevPC > 90:
		return rangeError("ActiveDevPC")
	case s.Temperature1 > 99:
		return rangeError("Temperature1")
	case s.Temperature2 > 99:
		return rangeError("Temperature2")
	case s.Temperature3 > 99:
		return rangeError("Temperature3")
	}
	return nil
}

// rangeError returns an out-of-range error.
func rangeError(source string) error {
	return errors.Errorf("%s out of range", source)
}

// Validate implements Validator.
func (v VisionStatus) Validate() error {
	switch v {
	case VisionStatusManual, VisionStatusOn, VisionStatusOff:
	default:
		return errors.Errorf("invalid VisionStatus: %s", v)
	}
	return nil
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
func (v CapacitorState) Validate() error {
	switch v {
	case CapacitorStateEmpty, CapacitorStateFull, CapacitorStateNoState:
	default:
		return errors.Errorf("invalid CapacitorState: %s", v)
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
