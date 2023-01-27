package common

import (
	"bytes"
	"io"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	body, err := buildPayload(config)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(strings.ToUpper(config.Method), url.String(), body)
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

func buildPayload(config HTTPConfig) (io.Reader, error) {
	var body io.Reader
	if config.Payload != nil {
		// req.Header.Add("content-type", "application/json")
		payload, err := json.Marshal(config.Payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(payload)
    }
	return body, nil
}
