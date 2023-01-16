package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(requestsCmd)
	requestsCmd.AddCommand(requestsListCmd)
}

var requestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "A set of commands to work with requests",
}

var requestsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all requests",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.Get("requests"))

		// for _, req := range viper.Get("requests") {
		// 	fmt.Println(req)
		// }
	},
}
