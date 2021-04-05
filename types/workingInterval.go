package types

import (
	"errors"
	"time"
)

type WorkingInterval struct {
	Start     time.Time
	End       time.Time
	WorkBreak time.Duration
}

func NewWorkingInterval(start time.Time, end time.Time, workBreak time.Duration) (*WorkingInterval, error) {
	if !end.IsZero() && start.After(end) {
		return &WorkingInterval{}, errors.New("End time must be after start time.")
	}

	return &WorkingInterval{start, end, workBreak}, nil
}
