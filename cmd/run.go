package cmd

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"net/http/httputil"
	"github.com/spf13/cobra"

	"github.com/michaeldbianchi/yahr/common"
)

func printRequest(req *http.Request) {
	reqDump, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: if !silent
	fmt.Printf("REQUEST:\n%s", string(reqDump))
}

func printResponse(execution common.RequestExecution) {
		// if !silent
		fmt.Println("Status:", execution.Response.StatusCode)
		// if !silent
		fmt.Println("Response Body:",execution.ResponseBody)
		// TODO: only print raw response if silent
		fmt.Println(execution.ResponseBody)
}

var runCmd = &cobra.Command{
	Use: "run",
	Short: "Execute http requests",
	Run: func(cmd *cobra.Command, arg []string) {
		requests := common.FetchRequestConfigs()
		config := requests[0]

		client := common.MakeClient(config)
		req, err := common.MakeRequest(config)
        if err != nil {
			log.Println("Failed to make request", err)
        }

		printRequest(req)

		execution, err := common.Execute(req, client)
		if err != nil {
			log.Fatal("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		printResponse(execution)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
