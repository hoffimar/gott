package cmd

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/hoffimar/gott/core"
	"github.com/hoffimar/gott/persistence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get statistics",
	Long:  `TODO.`,
	Run: func(cmd *cobra.Command, args []string) {
		var fileStore, _ = persistence.NewWorkingTimeFileStore(viper.GetString("StorageLocation"), "timerecording.json")
		var workingTimeList, _ = core.NewWorkingTimeList(fileStore)

		workingDayDuration, _ := time.ParseDuration("8h")

		statsPerDay, totalBalance, err := workingTimeList.GetWorkingTimeStatsPerDay(time.Now().Add(-statsSinceDuration), workingDayDuration)
		if err != nil {
			log.Fatal(err)
		}

		// create date slice to get sorted stats
		dates := make([]time.Time, len(statsPerDay))
		i := 0
		for k := range statsPerDay {
			dates[i] = k
			i++
		}

		sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

		for _, date := range dates {
			if statsPerDay[date].TotalBalance > 0 {
				color.Green("%s\t|Total=%s\t|Balance=%s\n", date, statsPerDay[date].Total.Round(time.Minute), statsPerDay[date].TotalBalance.Round(time.Minute))
			} else {
				color.Red("%s\t|Total=%s\t|Balance=%s\n", date, statsPerDay[date].Total.Round(time.Minute), statsPerDay[date].TotalBalance.Round(time.Minute))
			}
		}

		fmt.Println("============================================")
		fmt.Printf("Total Balance: %s\n", totalBalance.Round(time.Minute))
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	defaultDurationStatsSince, _ := time.ParseDuration("1w")
	statsCmd.Flags().DurationVar(&statsSinceDuration, "since", defaultDurationStatsSince, "Duration from when to evaluate the statistics.")
}
