package main

import (
	"errors"
	"fmt"
	"path"
)

// ErrInvalidHandler is thrown if any handler that is unknown to the program is specified
var ErrInvalidHandler = errors.New("configuration specifies unknown handler")

// ErrNoCommandsSpecified is thrown if the yaml file doesn't contain any commands
var ErrNoCommandsSpecified = errors.New("the specified yaml file doesn't contain any commands")

// ErrNoStrategiesSpecified is thrown if the yaml file doesn't contain any strategies
var ErrNoStrategiesSpecified = errors.New("the specified yaml file doesn't contain any strategies")

// ErrUndefinedCommand is thrown if a command specified can't be found in the yaml definition
var ErrUndefinedCommand = errors.New("the command you're trying to run doesn't exist in the yaml definition")

// ErrNoHandlerDefined is thrown when a strategy does not define a handler
var ErrNoHandlerDefined = errors.New("no handler specified in strategy")

// ErrInvalidStrategy is thrown when a invalid handler is referenced
var ErrInvalidStrategy = errors.New("invalid sttrategy specified in command")

var handlerFactories = map[string]func(map[string]interface{}) (Handler, error){
	"docker_run":          InitDockerRunHandler,
	"docker_compose_run":  InitComposeRunHandler,
	"docker_compose_exec": InitComposeExecHandler,
}

// Cfg is the uber object in our YAML file
type Cfg struct {
	commands       map[string]Handler
	handler        map[string]Handler
	defaultHandler Handler
}

// GenerateConfig is the main entry point from which we generate the config
func GenerateConfig(file []byte) (*Cfg, error) {
	yamlConfig, error := parseYaml(file)
	if error != nil {
		return nil, error
	}

	cfg := &Cfg{
		handler:  map[string]Handler{},
		commands: map[string]Handler{},
	}

	error = cfg.configFromYaml(yamlConfig)

	return cfg, error
}

// GetHandlerFor will try to find a handler for the specified command
func (cfg *Cfg) GetHandlerFor(command string, strictMode bool) (Handler, error) {
	executable := path.Base(command)
	if handler, ok := cfg.commands[executable]; ok {
		return handler, nil
	} else if strictMode {
		return nil, ErrUndefinedCommand
	} else if handler = cfg.defaultHandler; handler != nil {
		return handler, nil
	}
	return nil, ErrUndefinedCommand
}

// ListCommands allows for retrieval of all defined commands in a config
func (cfg *Cfg) ListCommands() []string {
	list := make([]string, 0, len(cfg.commands))
	for cmd := range cfg.commands {
		list = append(list, cmd)
	}
	return list
}

func (cfg *Cfg) configFromYaml(yaml *yamlCfg) error {
	if len(yaml.Strategies) == 0 {
		return ErrNoStrategiesSpecified
	}

	if len(yaml.Commands) == 0 {
		return ErrNoCommandsSpecified
	}

	for name, settings := range yaml.Strategies {
		if err := cfg.generateHandler(name, settings); err != nil {
			return err
		}
	}

	if name := yaml.DefaultStrategy; name != "" {
		if handler, ok := cfg.handler[name]; ok {
			cfg.defaultHandler = handler
		} else {
			return ErrInvalidHandler
		}
	}

	for command, strategy := range yaml.Commands {
		if handler, ok := cfg.handler[strategy]; ok {
			cfg.commands[command] = handler
		} else {
			return ErrInvalidStrategy
		}
	}

	return nil
}

func (cfg *Cfg) generateHandler(name string, settings map[string]interface{}) error {
	var handlerFactory func(map[string]interface{}) (Handler, error)
	if handlerName, ok := settings["handler"].(string); ok {
		handlerFactory, ok = handlerFactories[handlerName]
		if !ok {
			return ErrInvalidHandler
		}
	} else {
		return ErrNoHandlerDefined
	}

	delete(settings, "handler")

	if handler, err := handlerFactory(settings); err == nil {
		cfg.handler[name] = handler
	} else {
		return fmt.Errorf("Error in strategy %v: %v", name, err.Error())
	}

	return nil
}
