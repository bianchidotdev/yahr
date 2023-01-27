package common

import (
	"io"
	"log"
	"net/http"
)

type RequestExecution struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody string
}

// Run entire request based off of config
func Execute(reqConfig RequestConfig) (RequestExecution, error) {
	client := BuildClient(reqConfig.HTTPConfig)
	req, err := BuildHTTPRequest(reqConfig.HTTPConfig)
	if err != nil {
		return RequestExecution{Request: req}, err
	}

	log.Printf("Performing request for %s", reqConfig.Name)
	return PerformRequest(req, client)
}

func BuildClient(config HTTPConfig) *http.Client {
	client := &http.Client{}
	return client
}

func BuildHTTPRequest(config HTTPConfig) (*http.Request, error) {
	url := config.Url()
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	for name, value := range config.Headers {
		req.Header.Add(name, value)
	}

	return req, err
}

func PerformRequest(req *http.Request, client *http.Client) (RequestExecution, error) {
	execution := RequestExecution{
		Request: req,
	}
	execution.Request = req

	res, err := client.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
		return execution, err
	}
	execution.Response = res

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("client: could not read response body: %s\n", err)
		return execution, err
	}
	execution.ResponseBody = string(resBody)

	return execution, err
}
