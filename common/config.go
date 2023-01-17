package common

// import (
// 	"fmt"
// 	"gopkg.in/yaml.v2"
// 	"ioutil"
// 	"text/template"
// )

// type AppConfig struct {
// 	requests map[string]RequestConfig
// }

// type RequestConfig struct {
// 	method      string
// 	scheme      string
// 	host        string
// 	port        int
// 	path        string
// 	headers     map[string]string
// 	payload     []byte
// 	queryParams string
// }

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
// fmt.Println(str)										   //
// // unmarshal YAML here - from finalData.Bytes()		   //
/////////////////////////////////////////////////////////////
