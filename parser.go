package main

import (
	"errors"
	"gopkg.in/yaml.v2"
)

// ErrInvalidHandler is thrown if any handler that is unknown to the program is specified
var ErrInvalidHandler = errors.New("configuration specifies unknown handler")

// Current handler implementations
var availableHandlers = map[string]ExecHandler{"docker_compose_run": {"docker-compose", "run"}, "docker_compose_exec": {"docker-compose", "exec"}, "docker_run": {"docker", "run"}}

// ExecHandler is the desired OS exec
type ExecHandler struct {
	BaseCommand string
	Argument    string
}

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
	Image   string
}

// Validate checks whether a strategy specifies only valid handlers
func (s *Strategy) Validate() error {
	_, ok := availableHandlers[s.Handler]
	if !ok {
		return ErrInvalidHandler
	}

	return nil
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

	for _, strat := range cfg.Strategies {
		if err := strat.Validate(); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
