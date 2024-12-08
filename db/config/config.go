package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		DBName      string `yaml:"dbname"`
		MaxPoolSize int    `yaml:"max_pool_size"`
	} `yaml:"database"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("./db/config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
