package core

import (
	"errors"
	"time"

	"github.com/hoffimar/gott/types"
)

type WorkingTimeList struct {
	persistence WorkingTimeReadAdder
}

func NewWorkingTimeList(persistence WorkingTimeReadAdder) (list *WorkingTimeList, err error) {
	return &WorkingTimeList{persistence: persistence}, nil
}

func (list *WorkingTimeList) AddWorkingTimeInterval(interval types.WorkingInterval) (err error) {
	// TODO check if interval overlaps and return error

	return list.persistence.AddWorkingTime(interval)
}

func (list *WorkingTimeList) StartWorkingTimeInterval(startTime time.Time) (err error) {
	// TODO check if starting is possible, i.e. no other interval was already started

	defaultDuration, _ := time.ParseDuration("0m")
	interval, err := types.NewWorkingInterval(startTime, time.Time{}, defaultDuration)
	if err != nil {
		return err
	}
	return list.persistence.AddWorkingTime(*interval)
}

func (list *WorkingTimeList) StopWorkingTimeInterval(stopTime time.Time) (err error) {
	oldInterval, err := list.GetStartedWorkingTimeInterval()
	if err != nil {
		return err
	}

	newInterval, err := types.NewWorkingInterval(oldInterval.Start, stopTime, oldInterval.WorkBreak)
	if err != nil {
		return err
	}

	return list.persistence.UpdateWorkingTime(oldInterval, *newInterval)
}

func (list *WorkingTimeList) GetWorkingTimeIntervals() (workingTimes []types.WorkingInterval, err error) {
	return list.persistence.GetWorkingTimes()
}

func (list *WorkingTimeList) GetStartedWorkingTimeInterval() (interval types.WorkingInterval, err error) {
	times, err := list.persistence.GetWorkingTimes()
	if err != nil {
		return types.WorkingInterval{}, err
	}

	// Get interval without end time
	for idx := range times {
		element := &times[idx]
		if element.End.IsZero() {
			return *element, nil
		}
	}

	return types.WorkingInterval{}, errors.New("No started working time interval found.")
}
