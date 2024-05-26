package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type SensorConfig struct {
	ID                   int    `yaml:"id"`
	FilePath             string `yaml:"file_path"`
	TemperatureGaugeName string `yaml:"temperature_gauge_name"`
	HumidityGaugeName    string `yaml:"humidity_gauge_name"`
	QueryInterval        int    `yaml:"query_interval,omitempty"`
	ErrorInterval        int    `yaml:"error_interval,omitempty"`
}

type Config struct {
	Sensors []SensorConfig `yaml:"sensors"`
	Server  struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
}

var config Config

func loadConfig() {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		configFilePath = "config.yaml" // default path
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer configFile.Close()

	data, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
