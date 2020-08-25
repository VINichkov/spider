package environment

import (
	"github.com/dovadi/dbconfig"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)
type ConfigEnv map[string]map[string]string


type DbYamlConfig struct {
	Development map[string]string
	Production  map[string]string
}

//LoadYamlConfig is loading the yaml config file.
func LoadYamlConf(path string) *DbYamlConfig {
	var yamlconfig = DbYamlConfig{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &yamlconfig)
	if err != nil {
		panic(err)
	}

	return &yamlconfig
}



func LoadEnviroment(path string) {
	var environment string
	var env map[string]string
	jsonConf := dbconfig.LoadJSONConfig(path)
	conf := LoadYamlConf("config/application.yml")

	if len(jsonConf.Environment) == 0 {
		if len(os.Getenv("APPLICATION_ENV")) > 0 {
			environment = os.Getenv("APPLICATION_ENV")
		} else {
			environment = "development"
		}
	} else {
		environment = jsonConf.Environment
	}

	if environment == "development" {
		env = conf.Development
	} else {
		env = conf.Production
	}

	os.Setenv("APPLICATION_ENV", environment)

	for k, v := range env {
		os.Setenv(k, v)
	}
}