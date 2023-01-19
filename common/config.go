package common

import (
	"fmt"
	"log"
	"net/url"

	"github.com/imdario/mergo"
	"github.com/spf13/viper"
)

type Request struct {
	Name      string
	GroupName string
	RequestConfig
}

type RequestGroup struct {
	Name     string
	Requests []Request
	RequestConfig
}

type RequestConfig struct {
	Method      string
	Scheme      string
	Host        string
	Path        string
	Headers     map[string]string
	Payload     map[interface{}]interface{} //[]byte
	QueryParams string
}

func (config RequestConfig) Url() *url.URL {
	reqUrl := &url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   config.Path,
	}
	return reqUrl
}

func FetchRequestsByGroup(group string) []Request {
	groupConfig := FetchRequestGroup(group)

	return groupConfig.Requests
}

func FetchRequests() []Request {
	var requests []Request

	for _, group := range FetchRequestGroups() {
		for _, req := range group.Requests {
			requests = append(requests, req)
		}
	}

	return requests
}

func FetchRequestByName(group string, reqName string) Request {
	req, err := MakeRequest(group, reqName)
	if err != nil {
		log.Fatal("Failed to parse request", err)
	}

	return req
}

func FetchRequestGroups() []RequestGroup {
	var groups []RequestGroup

	for group, _ := range viper.GetStringMap("requests") {
		groupConfig := MakeRequestGroup(group)
		groups = append(groups, groupConfig)
	}

	return groups
}

func FetchRequestGroup(group string) RequestGroup {
	return MakeRequestGroup(group)
}

func MakeRequestGroup(group string) RequestGroup {
	var requests []Request
	groupAccessKey := fmt.Sprintf("requests.%s.requests", group)
	for key, _ := range viper.GetStringMap(groupAccessKey) {
		req, err := MakeRequest(group, key)
		if err != nil {
			log.Fatal("Failed to parse request", err)
		}

		requests = append(requests, req)
	}
	// TODO: handle group level configuration
	return RequestGroup{
		Name:     group,
		Requests: requests,
	}
}

func makeDefaultRequestConfig() RequestConfig {
	return RequestConfig{
		Method: "get",
		Scheme: "https",
		Path:   "/",
	}
}

func MakeRequest(group string, requestName string) (Request, error) {
	groupRequestConfig := makeDefaultRequestConfig()
	err := viper.UnmarshalKey(fmt.Sprintf("requests.%s", group), &groupRequestConfig)
	if err != nil {
		log.Println("Failed to parse group", err)
		return Request{}, err
	}

	accessKey := fmt.Sprintf("requests.%s.requests.%s", group, requestName)
	var requestConfig RequestConfig
	err = viper.UnmarshalKey(accessKey, &requestConfig)
	// TODO: this doesn't seem to be working
	if err != nil {
		log.Println("Failed to parse request", err)
		return Request{}, err
	}

	err = mergo.Merge(&requestConfig, groupRequestConfig)
	if err != nil {
		log.Fatal("Failed merging request configs", err)
	}

	request := Request{
		Name:          requestName,
		GroupName:     group,
		RequestConfig: requestConfig,
	}

	return request, nil
}

// func readConfig(path string) (ProxyConfig, error) {
// 	content, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		log.Fatal("Error when opening config file: ", err)
// 	}

// 	var config AppConfig
// 	err = yaml.Unmarshal(content, &config)
// 	if err != nil {
// 		log.Fatal("Error during Unmarshal(): ", err)
// 	}

// 	return config, err
// }

/////////////////////////////////////////////////////////////
// fileData, _ := ioutil.ReadFile("test.yml")
// var finalData bytes.Buffer							   //
// t := template.New("config")							   //
// t, err := t.Parse(string(fileData))					   //
// if err != nil {										   //
//     panic(err)										   //
// }													   //
// 														   //
// data := struct {										   //
//     THE_VARIABLE int									   //
// }{													   //
//     THE_VARIABLE: 30,  // replace with os.Getenv("FOO") //
// }													   //
// t.Execute(&finalData, data)							   //
// str := finalData.String()							   //
// log.Println(str)										   //
// // unmarshal YAML here - from finalData.Bytes()		   //
/////////////////////////////////////////////////////////////
