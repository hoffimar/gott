package core

import "github.com/hoffimar/gott/internal/types"

type WorkingTimeReadAdder interface {
	GetWorkingTimes() (intervals []types.WorkingInterval, err error)
	AddWorkingTime(interval types.WorkingInterval) (err error)
	UpdateWorkingTime(oldInterval types.WorkingInterval, newInterval types.WorkingInterval) (err error)
}
