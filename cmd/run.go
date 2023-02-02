package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/michaeldbianchi/yahr/core"
)

var RunCmd = &cli.Command{
	Name:      "run",
	Aliases:   []string{"r"},
	Usage:     "Execute HTTP requests",
	ArgsUsage: "GROUP REQUEST",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "silent", Aliases: []string{"s"}},
	},
	Action: func(cCtx *cli.Context) error {
		// TODO: implement select menu if not provided all options
		group := cCtx.Args().Get(0)
		name := cCtx.Args().Get(1)
		if group == "" || name == "" {
			fmt.Fprintf(cCtx.App.Writer, "No group or request specified\n\n")
			cli.ShowSubcommandHelp(cCtx)
			return nil
		}
		request := core.FetchRequestConfigByName(group, name)
		// TODO: how to handle errors effectively?
		// if err != nil {
		// 	log.Fatal("Could not find request", err)
		// }

		client := core.BuildClient(request.HTTPConfig)
		req, err := core.BuildHTTPRequest(request.HTTPConfig)
		if err != nil {
			log.Println("Failed to make request", err)
			return err
		}

		printRequest(cCtx, req)

		execution, err := core.PerformRequest(req, client)
		if err != nil {
			log.Println("client: could not create request:", err)
			return err
		}

		printResponse(cCtx, execution)
		return nil
	},
}

func printRequest(cCtx *cli.Context, req *http.Request) {
	reqDump, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		log.Fatal(err)
	}

	if !cCtx.Bool("silent") {
		fmt.Fprintf(cCtx.App.Writer, string(reqDump))
	}
}

func printResponse(cCtx *cli.Context, execution core.RequestExecution) {
	if !cCtx.Bool("silent") {
		fmt.Fprintf(cCtx.App.Writer, "Status: %d\n", execution.Response.StatusCode)
		fmt.Fprintf(cCtx.App.Writer, "Response Body:\n%s", execution.ResponseBody)
	} else {
		fmt.Fprintf(cCtx.App.Writer, execution.ResponseBody)
	}
}
