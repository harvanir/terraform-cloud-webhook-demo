package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const configName = "config.yml"

// Load load config from file
func Load() *AppConfig {
	appConfig := &AppConfig{}
	bytes, err := ioutil.ReadFile(configName)
	if err != nil {
		logrus.Fatal("failed when read file: %w", err)
	}
	// construct to type
	err = yaml.Unmarshal(bytes, appConfig)
	if err != nil {
		logrus.Fatal("failed when unmarshall: %w", err)
	}
	if appConfig.Terraform.Token.Bearer == "" {
		logrus.Fatal("no token found")
	}
	return appConfig
}
