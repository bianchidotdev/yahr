package cmd

import (
	"fmt"
	termtables "github.com/brettski/go-termtables"
	"github.com/spf13/cobra"

	"github.com/michaeldbianchi/yahr/common"
)

func printRequestList(requests []common.RequestConfig) {
	table := termtables.CreateTable()

	table.AddHeaders("Name", "Method", "Endpoint")
	for _, req := range requests {
		table.AddRow(req.Name, req.Method, req.Url())
	}
	fmt.Println(table.Render())
}

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
		requests := common.FetchRequestConfigs()
		printRequestList(requests)
	},
}
