package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const StatusSuccess string = "success"
const StatusFail string = "fail"

type Config struct {
	Server            ServerConfig             `yaml:"server"`
	BillingConnection AdapterConnectionAddress `yaml:"billing_connection"`
	LogLevel          string                   `yaml:"log_level"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type AdapterConnectionAddress struct {
	Address string `yaml:"address"`
}

var Service Config

func LoadConfig() error {
	// Read config from file
	yamlData, err := os.ReadFile("./config.yaml")
	if err != nil {
		return err
	}

	// Parse config data
	err = yaml.Unmarshal(yamlData, &Service)
	if err != nil {
		return err
	}

	return nil
}
