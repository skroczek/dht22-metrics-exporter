package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type EnvironmentData struct {
	Temperature float64 `json:"TempC"`
	Humidity    float64 `json:"Humidity"`
}

func readEnvironmentData(filePath string) (*EnvironmentData, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Check if the file is older than 5 minutes
	if time.Since(fileInfo.ModTime()) > 5*time.Minute {
		return nil, fmt.Errorf("data file %s is older than 5 minutes", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dataBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var data EnvironmentData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
