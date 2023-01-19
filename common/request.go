package common

import (
	"log"
	"io"
	"net/http"
)

type RequestExecution struct {
	Request *http.Request
	Response *http.Response
	ResponseBody string
}

func MakeClient(config RequestConfig) *http.Client {
	client := &http.Client{}
	return client
}

func MakeRequest(config RequestConfig) (*http.Request, error) {
	url := config.Url()
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	return req, err
}

func Execute(req *http.Request, client *http.Client) (RequestExecution, error) {
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

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("client: could not read response body: %s\n", err)
		return execution, err
	}
	execution.ResponseBody = string(resBody)

	return execution, err
}
