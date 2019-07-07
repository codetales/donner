package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// DockerRunHandler wraps a command with `docker container run`
type DockerRunHandler struct {
	Remove bool
	Image  string
}

// WrapCommand performs the actual wrapping of the command
func (handler *DockerRunHandler) WrapCommand(command []string) []string {
	wrappedCommand := []string{"docker", "container", "run"}

	if handler.Remove {
		wrappedCommand = append(wrappedCommand, "--rm")
	}

	if handler.Image != "" {
		wrappedCommand = append(wrappedCommand, handler.Image)
	}

	return append(wrappedCommand, command...)
}

func (handler *DockerRunHandler) validate() error {
	if handler.Image == "" {
		return errors.New("field image required but not set")
	}
	return nil
}

// InitDockerRunHandler generates a DockerRunHandler
func InitDockerRunHandler(settings map[string]interface{}) (Handler, error) {
	var handler *DockerRunHandler
	if err := mapstructure.Decode(settings, &handler); err != nil {
		return handler, err
	}

	if err := handler.validate(); err != nil {
		return handler, err
	}

	return handler, nil
}

// ComposeRunHandler wraps a command with `docker-compose exec`
type ComposeRunHandler struct {
	Remove  bool
	Service string
}

// WrapCommand performs the actual wrapping of the command
func (handler *ComposeRunHandler) WrapCommand(command []string) []string {
	wrappedCommand := []string{"docker-compose", "run"}

	if handler.Remove {
		wrappedCommand = append(wrappedCommand, "--rm")
	}

	if handler.Service != "" {
		wrappedCommand = append(wrappedCommand, handler.Service)
	}

	return append(wrappedCommand, command...)
}

func (handler *ComposeRunHandler) validate() error {
	if handler.Service == "" {
		return errors.New("field service required but not set")
	}
	return nil
}

// InitComposeRunHandler generates a ComposeRunHandler
func InitComposeRunHandler(settings map[string]interface{}) (Handler, error) {
	var handler *ComposeRunHandler
	var parsingMetadata *mapstructure.Metadata = &mapstructure.Metadata{}

	if err := mapstructure.DecodeMetadata(settings, &handler, parsingMetadata); err != nil {
		return handler, err
	}

	if err := handler.validate(); err != nil {
		return handler, err
	}

	if err := ensureNoAdditionalFields(parsingMetadata); err != nil {
		return handler, err
	}

	return handler, nil
}

func ensureNoAdditionalFields(m *mapstructure.Metadata) error {
	if len(m.Unused) > 0 {
		return fmt.Errorf("additonal field(s) detected: %v", strings.Join(m.Unused, ", "))
	}
	return nil
}

// ComposeExecHandler wraps a command with `docker-compose run`
type ComposeExecHandler struct {
	Service string
}

// WrapCommand performs the actual wrapping of the command
func (handler *ComposeExecHandler) WrapCommand(command []string) []string {
	wrappedCommand := []string{"docker-compose", "run"}

	if handler.Service != "" {
		wrappedCommand = append(wrappedCommand, handler.Service)
	}

	return append(wrappedCommand, command...)
}

func (handler *ComposeExecHandler) validate() error {
	if handler.Service == "" {
		return errors.New("field service required but not set")
	}
	return nil
}

// InitComposeExecHandler generates a ComposeExecHandler
func InitComposeExecHandler(settings map[string]interface{}) (Handler, error) {
	var handler *ComposeExecHandler
	if err := mapstructure.Decode(settings, &handler); err != nil {
		return handler, err
	}

	if err := handler.validate(); err != nil {
		return handler, err
	}

	return handler, nil
}
