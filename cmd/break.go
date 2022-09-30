package cmd

import (
	"log"
	"os"
	"time"

	"github.com/hoffimar/gott/internal/core"
	"github.com/hoffimar/gott/internal/persistence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "Add break time to the currently running interval.",
	Long:  `Only works when being checked in`,
	Run: func(cmd *cobra.Command, args []string) {

		var fileStore, _ = persistence.NewWorkingTimeFileStore(viper.GetString("StorageLocation"), "timerecording.json")
		var workingTimeList, _ = core.NewWorkingTimeList(fileStore)

		err := workingTimeList.AddBreakTime(breaktime)
		if err != nil {
			log.Fatal("Error adding break time: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)

	defaultDuration, _ := time.ParseDuration("5m")
	breakCmd.Flags().DurationVarP(&breaktime, "time", "t", defaultDuration, "the break time")
}
