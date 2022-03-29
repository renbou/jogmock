// Copyright 2022 Artem Mikheev

package activities

import "errors"

type ActivityType int

const (
	RunActivity ActivityType = iota + 1
	RideActivity
)

var ErrUnknownActivityType = errors.New("unknown activity type, only run and ride are known")

func (t ActivityType) IsKnown() bool {
	switch t {
	case RunActivity:
	case RideActivity:
	default:
		return false
	}
	return true
}
