package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/michaeldbianchi/yahr/common"
)

var RunCmd = &cli.Command{
	Name: "run",
	Aliases: []string{"r"},
	Usage: "Execute HTTP requests",
	ArgsUsage: "GROUP REQUEST",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "silent", Aliases: []string{"s"}},
	},
	Action: func(cCtx *cli.Context) error {
		// TODO: implement select menu if not provided all options
		group := cCtx.Args().Get(0)
		name := cCtx.Args().Get(1)
		request := common.FetchRequestConfigByName(group, name)
		// TODO: how to handle errors effectively?
		// if err != nil {
		// 	log.Fatal("Could not find request", err)
		// }

		client := common.BuildClient(request.HTTPConfig)
		req, err := common.BuildHTTPRequest(request.HTTPConfig)
		if err != nil {
			log.Println("Failed to make request", err)
			return err
		}

		printRequest(cCtx, req)

		execution, err := common.PerformRequest(req, client)
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
		fmt.Printf("Request:\n%s", string(reqDump))
	}
}

func printResponse(cCtx *cli.Context, execution common.RequestExecution) {
	if !cCtx.Bool("silent") {
		fmt.Println("Status:", execution.Response.StatusCode)
		fmt.Println("Response Body:\n", execution.ResponseBody)
	} else {
		fmt.Println(execution.ResponseBody)
	}
}

