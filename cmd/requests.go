package cmd

import (
	"fmt"
	termtables "github.com/brettski/go-termtables"
	"github.com/spf13/cobra"
	"log"

	"github.com/michaeldbianchi/yahr/common"
)

func printRequestList(requests []common.Request) {
	table := termtables.CreateTable()

	table.AddHeaders("Group", "Name", "Method", "Endpoint")
	for _, req := range requests {
		table.AddRow(req.GroupName, req.Name, req.Method, req.Url())
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
	Use:   "list [GROUP]",
	Short: "List all requests",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var requests []common.Request
		var err error
		if len(args) < 1 {
			requests = common.FetchRequests()
		} else {
			group := args[0]
			requests = common.FetchRequestsByGroup(group)
		}

		if err != nil {
			log.Fatal("Failed to load config", err)
		}
		printRequestList(requests)
	},
}
