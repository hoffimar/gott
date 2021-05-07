package core_test

import (
	"testing"
	"time"

	"github.com/hoffimar/gott/core"
	"github.com/hoffimar/gott/types"
)

type workingTimePersistenceMock struct {
	intervals []types.WorkingInterval
}

func (persistence *workingTimePersistenceMock) GetWorkingTimes() (intervals []types.WorkingInterval, err error) {
	return persistence.intervals, nil
}

func (persistence *workingTimePersistenceMock) AddWorkingTime(interval types.WorkingInterval) (err error) {
	persistence.intervals = append(persistence.intervals, interval)
	return nil
}

func (persistence *workingTimePersistenceMock) UpdateWorkingTime(oldInterval types.WorkingInterval, newInterval types.WorkingInterval) (err error) {
	return nil
}

func TestStats(t *testing.T) {

	currentTime := time.Now()
	sevenPM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 19, 0, 0, 0, time.Local)

	nowDate := time.Date(sevenPM.Year(), sevenPM.Month(), sevenPM.Day(), 0, 0, 0, 0, time.Local)

	interval1, _ := types.NewWorkingInterval(sevenPM.Add(-8*time.Hour), sevenPM.Add(-6*time.Hour), 0)
	intervalWithoutEnd, _ := types.NewWorkingInterval(sevenPM.Add(-5*time.Hour), time.Time{}, 0)

	tests := []struct {
		description         string
		intervals           []types.WorkingInterval
		since               time.Time
		targetDuration      time.Duration
		expectedStatsPerDay map[time.Time]*core.WorkingTimeStatsPerDay
		expectedBalance     time.Duration
		expectedErr         error
	}{
		{
			"Test negative balance",
			[]types.WorkingInterval{*interval1},
			time.Now().Add(-10 * time.Hour),
			8 * time.Hour,
			map[time.Time]*core.WorkingTimeStatsPerDay{
				nowDate: {Total: 2 * time.Hour, TotalBalance: -6 * time.Hour, StartTime: interval1.Start, EndTime: interval1.End},
			},
			-6 * time.Hour,
			nil,
		},
		{
			"Test running interval",
			[]types.WorkingInterval{*intervalWithoutEnd},
			time.Now().Add(-10 * time.Hour),
			8 * time.Hour,
			map[time.Time]*core.WorkingTimeStatsPerDay{
				nowDate: {Total: 5 * time.Hour, TotalBalance: -3 * time.Hour, StartTime: intervalWithoutEnd.Start, EndTime: intervalWithoutEnd.End},
			},
			-3 * time.Hour,
			nil,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.description, func(t *testing.T) {

			var persistenceMock = workingTimePersistenceMock{}
			persistenceMock.AddWorkingTime(testcase.intervals[0])
			var workingTimeList, _ = core.NewWorkingTimeList(&persistenceMock)

			statsPerDay, actualTotalBalance, actualError := workingTimeList.GetWorkingTimeStatsPerDay(testcase.since, sevenPM, testcase.targetDuration)
			if actualTotalBalance != testcase.expectedBalance {
				t.Errorf("Actual total balance %s != %s", actualTotalBalance, testcase.expectedBalance)
			}

			// check the individual stats per day
			actualStats, found := statsPerDay[nowDate]
			if !found {
				t.Errorf("No stats found for %s", nowDate)
			}

			if actualStats.Total != testcase.expectedStatsPerDay[nowDate].Total {
				t.Errorf("Total time %s, but expected %s", actualStats.Total, testcase.expectedStatsPerDay[nowDate].Total)
			}

			if actualStats.StartTime != testcase.expectedStatsPerDay[nowDate].StartTime {
				t.Errorf("Start time %s, but expected %s", actualStats.StartTime, testcase.expectedStatsPerDay[nowDate].StartTime)
			}

			if actualStats.EndTime != testcase.expectedStatsPerDay[nowDate].EndTime {
				t.Errorf("End time %s, but expected %s", actualStats.EndTime, testcase.expectedStatsPerDay[nowDate].EndTime)
			}

			if actualError != testcase.expectedErr {
				t.Errorf("Actual error %s, but expected %s", actualError, testcase.expectedErr)
			}
		})
	}
}
