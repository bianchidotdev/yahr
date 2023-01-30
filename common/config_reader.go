package common

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type MissingRequiredConfigError struct {
	Key string
}

func (e *MissingRequiredConfigError) Error() string {
	return fmt.Sprintf(`Missing required top-level key "%s"`, e.Key)
}

func ReadConfig(cfgFile string) ([]byte, error) {
	tmpl, err := template.ParseFiles(cfgFile)
	if err != nil {
		return nil, err
	}

	ydata := &bytes.Buffer{}

	environment, err := envToMap()
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(ydata, environment)
	if err != nil {
		return nil, err
	}

	return ydata.Bytes(), nil
}

func ParseConfig(configBytes []byte) (map[string]interface{}, error) {
	var appConfig map[string]interface{}
	err := yaml.Unmarshal(configBytes, &appConfig)
	if err != nil {
		return nil, err
	}

	return appConfig, nil
}

func SetConfig(appConfig map[string]interface{}) error {
	viper.Set("appConfig", appConfig)

	requestData := appConfig["requests"].(map[string]interface{})
	if requestData == nil {
		return &MissingRequiredConfigError{Key: "requests"}
	}
	viper.Set("requests", requestData)
	return nil
}

func envToMap() (map[string]string, error) {
	envMap := make(map[string]string)
	var err error

	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		envMap[split_v[0]] = strings.Join(split_v[1:], "=")
	}
	return envMap, err
}
