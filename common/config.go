package common

 import (
	"log"
	"net/url"
	"github.com/spf13/viper"
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

func (req RequestConfig) Url() *url.URL {
	reqUrl := &url.URL{
		Scheme: req.Scheme,
		Host:   req.Host,
		Path:   req.Path,
	}
	return reqUrl
}

func FetchRequestConfigs() []RequestConfig {
    var requests []RequestConfig

	for key, _ := range viper.GetStringMap("requests") {
		requests = append(requests, MakeRequestConfig(key))
	}

	return requests
}

func MakeRequestConfig(requestKey string) RequestConfig {
	request := RequestConfig{
		Name:   requestKey,
		Method: "get",
		Scheme: "https",
		Path:   "/",
	}
	accessKey := "requests." + requestKey

	err := viper.UnmarshalKey(accessKey, &request)
	if err != nil {
		log.Println("Failed to parse request", err)
	}
	return request
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
