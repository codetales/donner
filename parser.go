package main

import "gopkg.in/yaml.v2"

// Cfg is the uber object in our YAML file
type Cfg struct {
	Strategies      map[string]Strategy
	DefaultStrategy string `yaml:"default_strategy"`
	Commands        map[string]Command
}

// Strategy is the definition of a
type Strategy struct {
	Handler string
	Service string
	Remove  bool
}

// Command is an alias for string to properly reflect the yaml definition
type Command string

// ParseFile processes the .donner.yml file
func ParseFile(file []byte) (*Cfg, error) {
	cfg := Cfg{}
	err := yaml.Unmarshal([]byte(file), &cfg)

	if err != nil {
		return nil, err
	}

	// TODO validate config

	return &cfg, nil
}
