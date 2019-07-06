package main

import (
	"gopkg.in/yaml.v2"
)

type yamlCfg struct {
	Strategies      map[string]map[string]interface{}
	DefaultStrategy string `yaml:"default_strategy"`
	Commands        map[string]string
}

// parseFile processes the .donner.yml file
func parseYaml(file []byte) (*yamlCfg, error) {
	cfg := yamlCfg{}
	err := yaml.Unmarshal([]byte(file), &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
