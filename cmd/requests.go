package cmd

import (
	"fmt"
	termtables "github.com/brettski/go-termtables"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
)

type RequestConfig struct {
	Name        string
	Method      string
	Scheme      string
	Host        string
	Path        string
	Headers     map[string]string
	Payload     []byte
	QueryParams string
}

func (req RequestConfig) url() *url.URL {
	reqUrl := &url.URL{
		Scheme: req.Scheme,
		Host:   req.Host,
		Path:   req.Path,
	}
	return reqUrl
}

func newRequest(requestKey string) RequestConfig {
	request := RequestConfig{
		Name:   requestKey,
		Method: "get",
		Scheme: "https",
		Path:   "/",
	}
	accessKey := "requests." + requestKey

	err := viper.UnmarshalKey(accessKey, &request)
	if err != nil {
		fmt.Println("Failed to parse request", err)
	}
	return request
}

func printRequestList(requests []RequestConfig) {
	table := termtables.CreateTable()

	table.AddHeaders("Name", "Method", "Endpoint")
	for _, req := range requests {
		table.AddRow(req.Name, req.Method, req.url())
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
		var requests []RequestConfig
		for key, _ := range viper.GetStringMap("requests") {
			requests = append(requests, newRequest(key))
		}
		printRequestList(requests)
	},
}
