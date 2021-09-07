/*
	Copyright 2021 Misir Jafarov

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

			http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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
