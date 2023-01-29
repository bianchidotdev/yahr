package cmd

import (
	"fmt"
	termtables "github.com/brettski/go-termtables"
	"github.com/urfave/cli/v2"

	"github.com/michaeldbianchi/yahr/common"
)

func printRequestList(requests []common.RequestConfig) {
	table := termtables.CreateTable()

	table.AddHeaders("Group", "Name", "Method", "Endpoint")
	for _, req := range requests {
		table.AddRow(req.GroupName, req.Name, req.Method, req.Url())
	}
	fmt.Println(table.Render())
}

var RequestCmd = &cli.Command{
	Name: "requests",
	Aliases: []string{"req",},
	Usage: "A set of commands to work with requests",
	Subcommands: []*cli.Command{
		requestListCmd,
	},
}

var requestListCmd = &cli.Command{
	Name: "list",
	Aliases: []string{"l",},
	Usage: "List all requests, optionally limited to group of requests",
	ArgsUsage: "[GROUP]",
	Action: func(cCtx *cli.Context) error {
		var requests []common.RequestConfig
		var err error
		if cCtx.NArg() < 1 {
			requests = common.FetchRequestConfigs()
		} else {
			group := cCtx.Args().First()
			requests = common.FetchRequestConfigsByGroup(group)
		}

		if err != nil {
			return err
		}
		printRequestList(requests)
		return nil
	},
}
