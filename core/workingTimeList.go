package core

import (
	"errors"
	"sort"
	"time"

	"github.com/hoffimar/gott/types"
)

var ErrNoIntervalStarted = errors.New("No working interval started.")

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

	return types.WorkingInterval{}, ErrNoIntervalStarted
}

type WorkingTimeStatsPerDay struct {
	Total        time.Duration
	TotalBalance time.Duration
	StartTime    time.Time
	EndTime      time.Time
	// TODO add break duration
}

func (list *WorkingTimeList) GetWorkingTimeStatsPerDay(since time.Time, until time.Time, targetWorkingTime time.Duration) (result map[time.Time]*WorkingTimeStatsPerDay, totalBalance time.Duration, err error) {
	times, err := list.GetWorkingTimeIntervals()
	if err != nil {
		return nil, 0, err
	}

	result = make(map[time.Time]*WorkingTimeStatsPerDay)

	sort.Slice(times, func(i, j int) bool { return times[i].Start.Before(times[j].Start) })

	for _, interval := range times {
		// TODO filter working times here, should be moved to storage
		// TODO filter times after 'until', limit the balance
		if !interval.Start.Before(since) {

			year, month, day := interval.Start.Date()
			date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			var total time.Duration = 0
			if interval.End.IsZero() {
				total = until.Sub(interval.Start) - interval.WorkBreak
			} else {
				total = interval.End.Sub(interval.Start) - interval.WorkBreak
			}

			element, found := result[date]
			if found {
				element.Total = element.Total + total
				element.TotalBalance = element.TotalBalance + total
				if interval.Start.Before(element.StartTime) {
					element.StartTime = interval.Start
				}
				if interval.End.After(element.EndTime) {
					element.EndTime = interval.End
				}
			} else {
				result[date] = &WorkingTimeStatsPerDay{Total: total, TotalBalance: total - targetWorkingTime, StartTime: interval.Start, EndTime: interval.End}
			}
		}
	}

	// Calculate total balance
	for _, v := range result {
		totalBalance += v.TotalBalance
	}

	return result, totalBalance, nil
}
