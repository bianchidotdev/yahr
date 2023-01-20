package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/michaeldbianchi/yahr/common"
)

func printRequest(req *http.Request) {
	reqDump, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		log.Fatal(err)
	}

	if !viper.GetBool("silent") {
		fmt.Printf("Request:\n%s", string(reqDump))
	}
}

func printResponse(execution common.RequestExecution) {
	if !viper.GetBool("silent") {
		fmt.Println("Status:", execution.Response.StatusCode)
		fmt.Println("Response Body:\n", execution.ResponseBody)
	} else {
		fmt.Println(execution.ResponseBody)
	}
}

var runCmd = &cobra.Command{
	Use:   "run GROUP REQUEST",
	Short: "Execute http requests",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, arg []string) {
		// TODO: implement select menu if not provided all options
		group := arg[0]
		name := arg[1]
		request := common.FetchRequestConfigByName(group, name)
		// TODO: how to handle errors effectively?
		// if err != nil {
		// 	log.Fatal("Could not find request", err)
		// }

		client := common.MakeClient(request.HTTPConfig)
		req, err := common.MakeHTTPRequest(request.HTTPConfig)
		if err != nil {
			log.Fatal("Failed to make request", err)
		}

		printRequest(req)

		execution, err := common.Execute(req, client)
		if err != nil {
			log.Fatal("client: could not create request:", err)
			os.Exit(1)
		}

		printResponse(execution)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
