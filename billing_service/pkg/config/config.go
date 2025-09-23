package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database   DatabaseConfig   `yaml:"database"`
	GRPCServer GRPCServerConfig `yaml:"grpc_server"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"db_name"`
}

type GRPCServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var Service Config

func LoadConfig() error {
	// read config from file
	yamlData, err := os.ReadFile("./config.yaml")
	if err != nil {
		return err
	}

	// pass config data to BillingConfig
	err = yaml.Unmarshal(yamlData, &Service)
	return err
}
