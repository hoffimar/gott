package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"time"

	"github.com/hoffimar/gott/persistence"
	"github.com/hoffimar/gott/types"
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
		fmt.Println("workingtime add called")

		os.MkdirAll(viper.GetString("StorageLocation"), 0700)

		var file *os.File
		var err error
		filePath := path.Join(viper.GetString("StorageLocation"), "timerecording.json")
		// Check for file existence
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Creating file %s", filePath)

			// fill file with initial array
			content := []byte("[]")
			ioutil.WriteFile(filePath, content, 0600)
			if err != nil {
				log.Fatal(err)
			}
		}

		file, _ = os.OpenFile(filePath, os.O_RDWR, 0600)

		// first get existing time recordings, then add the new one
		times, err := persistence.GetWorkingTimes(path.Join(viper.GetString("StorageLocation"), "timerecording.json"))
		if err != nil {
			fmt.Println("Could not read existing times: ", err)
		}

		// parse times
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

		// TODO check that no existing time overlaps, present a warning
		//inputInterval := types.WorkingInterval{Start: startTime, End: endTime, WorkBreak: breaktime}
		inputInterval, err := types.NewWorkingInterval(startTime, endTime, breaktime)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		times = append(times, *inputInterval)

		defer file.Close()
		persistence.SaveWorkingTimes(file, times)
	},
}

func getTimeFromInputString(input string) (result time.Time, err error) {
	result, err = time.Parse("2006-01-02-15:04", input)
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
		fmt.Println("workingtime log called")

		times, err := persistence.GetWorkingTimes(path.Join(viper.GetString("StorageLocation"), "timerecording.json"))
		if err != nil {
			fmt.Println("Error reading times: ", err)
		}

		sort.Slice(times, func(i, j int) bool { return times[i].Start.Before(times[j].Start) })

		for _, interval := range times {
			fmt.Printf("Start: %s, end: %s, break: %s\n", interval.Start, interval.End, interval.WorkBreak)
		}
	},
}

var statsWorkingTimeCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get working time statistics",
	Long:  `TODO.`,
	Run: func(cmd *cobra.Command, args []string) {
		times, err := persistence.GetWorkingTimes(path.Join(viper.GetString("StorageLocation"), "timerecording.json"))
		if err != nil {
			fmt.Println("Error reading times: ", err)
		}

		sort.Slice(times, func(i, j int) bool { return times[i].Start.Before(times[j].Start) })

		var totalBalance time.Duration
		for _, interval := range times {
			// filter working times here, should be moved to storage
			if time.Since(interval.Start) < statsSinceDuration {
				workingDayDuration, _ := time.ParseDuration("8h")
				balance := interval.End.Sub(interval.Start) - interval.WorkBreak - workingDayDuration
				totalBalance += balance
				fmt.Printf("Start: %s, end: %s, break: %s, --- balance: %s\n", interval.Start, interval.End, interval.WorkBreak, balance)
			}
		}
		fmt.Printf("TOTAL BALANCE: %s", totalBalance)
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
	workingtimeCmd.AddCommand(statsWorkingTimeCmd)

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

	defaultDurationStatsSince, _ := time.ParseDuration("1w")
	statsWorkingTimeCmd.Flags().DurationVar(&statsSinceDuration, "since", defaultDurationStatsSince, "Duration from when to evaluate the statistics.")
}
