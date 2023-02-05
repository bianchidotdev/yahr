package core

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

type RequestNotFoundError struct {
	Path string
}

func (m *RequestNotFoundError) Error() string {
	return fmt.Sprintf("Cannot find path '%s' in yaml config", m.Path)
}

type InvalidHTTPMethodError struct {
	Method string
}

func (m *InvalidHTTPMethodError) Error() string {
	return fmt.Sprintf("Method %s not valid", strings.ToUpper(m.Method))
}

type RequestConfig struct {
	Name      string
	GroupName string
	*HTTPConfig
}

type RequestGroup struct {
	Name     string
	Requests []*RequestConfig
	*HTTPConfig
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

func (config HTTPConfig) InterpolatePathParams() (string, error) {
	path := config.Path
	regex := regexp.MustCompile(`/:(\w+)`)
	matches := regex.FindAllStringSubmatch(path, -1)
	if matches == nil {
		return path, nil
	}

	params := viper.GetStringMapString("pathParams")
	for _, match := range matches {
		match_string := match[1]
		if match_string == "" {
			return "", fmt.Errorf("empty path param: %s", match[0])
		}
		replacement, ok := params[match_string]
		if !ok {
			// TODO: maybe don't error out here and delegate errors to the validatiom method
			return "", fmt.Errorf("missing required path param '%s' - specify with '-p %s=<value>'", match_string, match_string)
		}
		path = strings.ReplaceAll(path, fmt.Sprintf(":%s", match_string), replacement)
	}

	return path, nil
}

func FetchRequestConfigsByGroup(group string) []*RequestConfig {
	groupConfig := fetchRequestGroup(group)

	return groupConfig.Requests
}

func FetchRequestConfigs() []*RequestConfig {
	var requests []*RequestConfig

	for _, group := range fetchRequestGroups() {
		for _, req := range group.Requests {
			requests = append(requests, req)
		}
	}

	return requests
}

func FetchRequestConfigByName(group string, reqName string) (*RequestConfig, error) {
	req, err := makeRequestConfig(group, reqName)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func fetchRequestGroups() []*RequestGroup {
	var groups []*RequestGroup

	for group, _ := range viper.GetStringMap("requests") {
		groupConfig := makeRequestGroup(group)
		groups = append(groups, groupConfig)
	}

	return groups
}

func fetchRequestGroup(group string) *RequestGroup {
	return makeRequestGroup(group)
}

func makeRequestGroup(group string) *RequestGroup {
	var requests []*RequestConfig
	groupAccessKey := fmt.Sprintf("requests.%s.requests", group)
	for key, _ := range viper.GetStringMap(groupAccessKey) {
		req, err := makeRequestConfig(group, key)
		if err != nil {
			log.Fatal("Failed to parse request", err)
		}

		requests = append(requests, req)
	}
	return &RequestGroup{
		Name:     group,
		Requests: requests,
	}
}

func makeDefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Method: "get",
		Scheme: "https",
		Path:   "/",
	}
}

func makeRequestConfig(group string, requestName string) (*RequestConfig, error) {
	groupAccessKey := fmt.Sprintf("requests.%s", group)
	if !viper.IsSet(groupAccessKey) {
		return nil, &RequestNotFoundError{Path: groupAccessKey}
	}

	groupHTTPConfig := makeDefaultHTTPConfig()
	err := viper.UnmarshalKey(groupAccessKey, &groupHTTPConfig)
	if err != nil {
		log.Println("Failed to parse group", err)
		return nil, err
	}

	accessKey := fmt.Sprintf("requests.%s.requests.%s", group, requestName)
	if !viper.IsSet(accessKey) {
		return nil, &RequestNotFoundError{Path: accessKey}
	}

	var httpConfig *HTTPConfig
	err = viper.UnmarshalKey(accessKey, &httpConfig)
	if err != nil {
		return nil, err
	}

	err = mergo.Merge(httpConfig, groupHTTPConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed merging request configs: %s", err)
	}

	httpConfig.Path, err = httpConfig.InterpolatePathParams()
	if err != nil {
		return nil, err
	}

	err = validateHTTPConfig(httpConfig)
	if err != nil {
		return nil, err
	}

	request := &RequestConfig{
		Name:       requestName,
		GroupName:  group,
		HTTPConfig: httpConfig,
	}

	return request, nil
}

// TODO: make function on struct
func validateHTTPConfig(config *HTTPConfig) error {
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
