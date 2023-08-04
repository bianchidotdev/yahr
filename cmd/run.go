package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/bianchidotdev/yahr/core"
)

func setPathParams(cCtx *cli.Context) error {
	var pathParams = make(map[string]string)
	for _, param := range cCtx.StringSlice("path-param") {
		split := strings.Split(param, "=")
		if split[0] == "" || split[1:] == nil {
			return fmt.Errorf("Invalid path params - must be specified in key=value pairs")
		}
		value := strings.Join(split[1:], "=")
		pathParams[split[0]] = value
	}
	viper.Set("pathParams", pathParams)
	return nil
}

var RunCmd = &cli.Command{
	Name:      "run",
	Aliases:   []string{"r"},
	Usage:     "Execute HTTP requests",
	ArgsUsage: "GROUP REQUEST",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "silent", Aliases: []string{"s"}},
		&cli.StringSliceFlag{
			Name:    "path-param",
			Usage:   "params for path interpolation - key=value",
			Aliases: []string{"p"},
		},
	},
	Action: func(cCtx *cli.Context) error {
		group := cCtx.Args().Get(0)
		name := cCtx.Args().Get(1)
		if group == "" || name == "" {
			// TODO: implement select menu if not provided all options
			fmt.Fprintf(cCtx.App.Writer, "No group or request specified\n\n")
			cli.ShowSubcommandHelp(cCtx)
			return nil
		}
		err := setPathParams(cCtx)
		if err != nil {
			return err
		}

		request, err := core.FetchRequestConfigByName(group, name)
		if err != nil {
			return err
		}

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
		if execution.Response.StatusCode >= 400 {
			return cli.Exit("", 1)
		}
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

func printResponse(cCtx *cli.Context, execution *core.RequestExecution) {
	if !cCtx.Bool("silent") {
		fmt.Fprintf(cCtx.App.Writer, "Status: %d\n", execution.Response.StatusCode)
		fmt.Fprintf(cCtx.App.Writer, "Response Body:\n%s", execution.ResponseBody)
	} else {
		fmt.Fprintf(cCtx.App.Writer, execution.ResponseBody)
	}
}
