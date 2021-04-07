package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/hoffimar/gott/core"
	"github.com/hoffimar/gott/persistence"
	"github.com/hoffimar/gott/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of today's working time.",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		var fileStore, _ = persistence.NewWorkingTimeFileStore(viper.GetString("StorageLocation"), "timerecording.json")
		var workingTimeList, _ = core.NewWorkingTimeList(fileStore)

		now := time.Now().Round(time.Second)
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		workingDayDuration, _ := time.ParseDuration("8h")

		startedInterval, err := workingTimeList.GetStartedWorkingTimeInterval()
		if errors.Is(err, core.ErrNoIntervalStarted) {
			fmt.Printf("\nNo working interval started.\n\n")
		} else {
			if err != nil {
				fmt.Printf("Error determining whether interval is running: %s", err)
				os.Exit(1)
			}

			color.Blue("\nCurrent work time started at %s, %s ago\n\n", startedInterval.Start, now.Sub(startedInterval.Start))
		}

		// Get today's working time so far
		times, err := workingTimeList.GetWorkingTimeIntervals()
		if err != nil {
			fmt.Println("Error reading times: ", err)
		}

		var todaysTimes []types.WorkingInterval
		for i := range times {
			if times[i].Start.After(today) {
				todaysTimes = append(todaysTimes, times[i])
			}
		}

		var total time.Duration
		for _, interval := range todaysTimes {
			if interval.End.IsZero() {
				total += now.Sub(interval.Start) - interval.WorkBreak
			} else {
				total += interval.End.Sub(interval.Start) - interval.WorkBreak
			}
		}

		balance := total - workingDayDuration
		durationZero, _ := time.ParseDuration("0m")
		if balance < durationZero {
			fmt.Printf("Balance today: %s (total so far: %s), estimated end time: %s\n", balance, total, now.Add(-balance))
		} else {
			color.Green("Balance today: %s (total so far: %s)\n", balance, total)
		}

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
