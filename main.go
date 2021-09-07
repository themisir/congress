package main

import (
	"congress/logger"
	"congress/models"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func readConfig(filename string) (*models.Config, error) {
	config := new(models.Config)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	config, err := readConfig("congress.yaml")
	if err != nil {
		logger.Default.Error("Failed to read file: %s", err)
	} else {
		app := &App{config}
		app.Run()
	}
}
