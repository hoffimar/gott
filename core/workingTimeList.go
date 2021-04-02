package core

import (
	"github.com/hoffimar/gott/persistence"
	"github.com/hoffimar/gott/types"
)

type WorkingTimeList struct {
	persistence persistence.WorkingTimeReadAdder
}

func NewWorkingTimeList(persistence persistence.WorkingTimeReadAdder) (list *WorkingTimeList, err error) {
	return &WorkingTimeList{persistence: persistence}, nil
}

func (list *WorkingTimeList) AddWorkingTimeInterval(interval types.WorkingInterval) (err error) {
	// TODO check if interval overlaps and return error

	return list.persistence.AddWorkingTime(interval)
}

func (list *WorkingTimeList) GetWorkingTimeIntervals() (workingTimes []types.WorkingInterval, err error) {
	return list.persistence.GetWorkingTimes()
}
