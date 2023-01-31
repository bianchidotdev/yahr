package core

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/imdario/mergo"
	"github.com/spf13/viper"
)

type InvalidHTTPMethodError struct {
	Method string
}

func (m *InvalidHTTPMethodError) Error() string {
	return fmt.Sprintf("Method %s not valid", strings.ToUpper(m.Method))
}

type RequestConfig struct {
	Name      string
	GroupName string
	HTTPConfig
}

type RequestGroup struct {
	Name     string
	Requests []RequestConfig
	HTTPConfig
}

type HTTPConfig struct {
	Method      string
	Scheme      string
	Host        string
	Path        string
	Headers     map[string]string
	Payload     map[string]interface{} //[]byte
	QueryParams map[string]string      `yaml:"query_params"`
}

func (config HTTPConfig) Url() *url.URL {
	reqUrl := &url.URL{
		Scheme: config.Scheme,
		Host:   config.Host,
		Path:   config.Path,
	}

	if config.QueryParams != nil {
		query := url.Values{}
		for key, value := range config.QueryParams {
			query.Add(key, value)
		}
		reqUrl.RawQuery = query.Encode()
	}

	return reqUrl
}

func FetchRequestConfigsByGroup(group string) []RequestConfig {
	groupConfig := fetchRequestGroup(group)

	return groupConfig.Requests
}

func FetchRequestConfigs() []RequestConfig {
	var requests []RequestConfig

	for _, group := range fetchRequestGroups() {
		for _, req := range group.Requests {
			requests = append(requests, req)
		}
	}

	return requests
}

func FetchRequestConfigByName(group string, reqName string) RequestConfig {
	req, err := makeRequestConfig(group, reqName)
	if err != nil {
		log.Fatal("Failed to parse request - ", err)
	}

	return req
}

func fetchRequestGroups() []RequestGroup {
	var groups []RequestGroup

	for group, _ := range viper.GetStringMap("requests") {
		groupConfig := makeRequestGroup(group)
		groups = append(groups, groupConfig)
	}

	return groups
}

func fetchRequestGroup(group string) RequestGroup {
	return makeRequestGroup(group)
}

func makeRequestGroup(group string) RequestGroup {
	var requests []RequestConfig
	groupAccessKey := fmt.Sprintf("requests.%s.requests", group)
	for key, _ := range viper.GetStringMap(groupAccessKey) {
		req, err := makeRequestConfig(group, key)
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

func makeDefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Method: "get",
		Scheme: "https",
		Path:   "/",
	}
}

func makeRequestConfig(group string, requestName string) (RequestConfig, error) {
	groupHTTPConfig := makeDefaultHTTPConfig()
	err := viper.UnmarshalKey(fmt.Sprintf("requests.%s", group), &groupHTTPConfig)
	if err != nil {
		log.Println("Failed to parse group", err)
		return RequestConfig{}, err
	}

	accessKey := fmt.Sprintf("requests.%s.requests.%s", group, requestName)
	var httpConfig HTTPConfig
	err = viper.UnmarshalKey(accessKey, &httpConfig)
	// TODO: this doesn't seem to be working
	if err != nil {
		log.Println("Failed to parse request", err)
		return RequestConfig{}, err
	}

	err = mergo.Merge(&httpConfig, groupHTTPConfig)
	if err != nil {
		log.Fatal("Failed merging request configs", err)
	}

	err = validateHTTPConfig(httpConfig)
	if err != nil {
		return RequestConfig{}, err
	}

	request := RequestConfig{
		Name:       requestName,
		GroupName:  group,
		HTTPConfig: httpConfig,
	}

	return request, nil
}

func validateHTTPConfig(config HTTPConfig) error {
	httpMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	// TODO: include debug details to make it easy to fix
	if !slices.Contains(httpMethods, strings.ToUpper(config.Method)) {
		return &InvalidHTTPMethodError{Method: config.Method}
	}

	return nil
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
