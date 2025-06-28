package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var AppCfg AppConfig

func (this *AppConfig) ReadConfigFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to read config file: %v", err))
	}
	if err := yaml.Unmarshal(data, this); err != nil {
		panic(fmt.Sprintf("failed to parse yaml: %v", err))
	}
}
