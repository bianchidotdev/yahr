package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	Use: "run REQUEST",
	Short: "Execute http requests",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, arg []string) {
		config, err := common.FetchRequestConfigByName(arg[0])
		if err != nil {
			log.Fatal("Could not find request", err)
		}

		client := common.MakeClient(config)
		req, err := common.MakeRequest(config)
		if err != nil {
			log.Fatal("Failed to make request", err)
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
