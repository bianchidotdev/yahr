package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of yahr",
	Long:  `Print the version information of yahr`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: make version not dependent on config file
		fmt.Println("yahr v0.0.1 -- HEAD")
	},
}
