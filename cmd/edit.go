package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit working times by opening an editor with the configured storage file.",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {

		filePath := path.Join(viper.GetString("StorageLocation"), "timerecording.json")
		editor := viper.GetString("Editor")

		command := exec.Command(editor, filePath)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err := command.Run()
		if err != nil {
			fmt.Printf("Error opening editor %s with file %s", editor, filePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
