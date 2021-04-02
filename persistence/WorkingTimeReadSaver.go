package persistence

import "github.com/hoffimar/gott/types"

type WorkingTimeReadAdder interface {
	GetWorkingTimes() (intervals []types.WorkingInterval, err error)
	AddWorkingTime(interval types.WorkingInterval) (err error)
}
