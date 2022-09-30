package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/hoffimar/gott/internal/core"
	"github.com/hoffimar/gott/internal/persistence"
	"github.com/hoffimar/gott/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// workingtimeCmd represents the workingtime command
var workingtimeCmd = &cobra.Command{
	Use:   "workingtime",
	Short: "A brief description of your command",
	Long:  `A longer description `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("workingtime called")
	},
	TraverseChildren: true,
}

// addCmd represents the add command
var addWorkingTimeCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a working time (with start and end time) to the file.",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		// parse input times
		startTime, err := getTimeFromInputString(startTimeString)
		if err != nil {
			fmt.Printf("start time parsing not possible: %s", err)
			os.Exit(1)
		}

		endTime, err := getTimeFromInputString(endTimeString)
		if err != nil {
			fmt.Printf("end time parsing not possible: %s", err)
			os.Exit(1)
		}

		inputInterval, err := types.NewWorkingInterval(startTime, endTime, breaktime)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		var fileStore, _ = persistence.NewWorkingTimeFileStore(viper.GetString("StorageLocation"), "timerecording.json")
		var workingTimeList, _ = core.NewWorkingTimeList(fileStore)

		err = workingTimeList.AddWorkingTimeInterval(*inputInterval)
		if err != nil {
			log.Fatal("Error adding the working time: ", err)
		}
	},
}

func getTimeFromInputString(input string) (result time.Time, err error) {
	result, err = time.ParseInLocation("2006-01-02-15:04", input, time.Local)
	if err != nil {
		var inputTime time.Time
		inputTime, err = time.Parse("15:04", input)
		if err != nil {
			fmt.Printf("time parsing not possible for %s", input)
			return time.Time{}, err // return empty time struct
		}

		// set date to today
		today := time.Now()
		result = time.Date(today.Year(), today.Month(), today.Day(), inputTime.Hour(), inputTime.Minute(), 0, 0, time.Local)
	}

	return result, nil
}

var logWorkingTimesCmd = &cobra.Command{
	Use:   "log",
	Short: "Show a log of the times entered, sorted by date",
	Long:  `TODO.`,
	Run: func(cmd *cobra.Command, args []string) {
		var fileStore, _ = persistence.NewWorkingTimeFileStore(viper.GetString("StorageLocation"), "timerecording.json")
		var workingTimeList, _ = core.NewWorkingTimeList(fileStore)

		times, err := workingTimeList.GetWorkingTimeIntervals()
		if err != nil {
			fmt.Println("Error reading times: ", err)
		}

		sort.Slice(times, func(i, j int) bool { return times[i].Start.Before(times[j].Start) })

		for _, interval := range times {
			fmt.Printf("Start: %s, end: %s, break: %s\n", interval.Start, interval.End, interval.WorkBreak)
		}
	},
}

var (
	breaktime          time.Duration
	startTimeString    string
	endTimeString      string
	statsSinceDuration time.Duration
)

func init() {

	rootCmd.AddCommand(workingtimeCmd)
	workingtimeCmd.AddCommand(addWorkingTimeCmd)
	workingtimeCmd.AddCommand(logWorkingTimesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workingtimeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	defaultDuration, _ := time.ParseDuration("0m")
	addWorkingTimeCmd.Flags().StringVarP(&startTimeString, "start", "s", "", "Start time of the working time interval")
	addWorkingTimeCmd.Flags().StringVarP(&endTimeString, "end", "e", "", "End time of the working time interval")
	addWorkingTimeCmd.Flags().DurationVarP(&breaktime, "break", "b", defaultDuration, "the break time")
}
