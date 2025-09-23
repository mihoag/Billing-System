package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const StatusSuccess string = "success"
const StatusFail string = "fail"

type Config struct {
	Logger            *zap.Logger              `yaml:"-"`
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
	// Get config file path from environment or use default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	// Read config from file
	yamlData, err := os.ReadFile(configPath)
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
