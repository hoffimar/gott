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

var checkInCmd = &cobra.Command{
	Use:   "checkin",
	Short: "Check-in, i.e., start a new working time interval using the current time as starting time.",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		var fileStore, _ = persistence.NewWorkingTimeFileStore(viper.GetString("StorageLocation"), "timerecording.json")
		var workingTimeList, _ = core.NewWorkingTimeList(fileStore)

		err := workingTimeList.StartWorkingTimeInterval(time.Now().Round(time.Second))
		if err != nil {
			log.Fatal("Error starting the working time: ", err)
			os.Exit(1)
		}
	},
}

var checkOutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Check-out, i.e., stop the current working time interval.",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		var fileStore, _ = persistence.NewWorkingTimeFileStore(viper.GetString("StorageLocation"), "timerecording.json")
		var workingTimeList, _ = core.NewWorkingTimeList(fileStore)

		err := workingTimeList.StopWorkingTimeInterval(time.Now().Round(time.Second))
		if err != nil {
			log.Fatal("Error stopping the working time: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkInCmd)
	rootCmd.AddCommand(checkOutCmd)
}
