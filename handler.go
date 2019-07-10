package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Handler is the interface around the individual handler implementations
type Handler interface {
	BuildCommand([]string) []string
	validate() error
}

// FallbackHandler does not wrap the command and will just execute it as is
type FallbackHandler struct{}

// BuildCommand simply returns the passed in command
func (h *FallbackHandler) BuildCommand(command []string) []string {
	return command
}

func (h *FallbackHandler) validate() error {
	return nil
}

// DockerRunHandler wraps a command with `docker container run`
type DockerRunHandler struct {
	Remove bool
	Image  string
}

// BuildCommand performs the actual wrapping of the command
func (h *DockerRunHandler) BuildCommand(command []string) []string {
	wrappedCommand := []string{"docker", "run", "-it"}

	if h.Remove {
		wrappedCommand = append(wrappedCommand, "--rm")
	}

	if h.Image != "" {
		wrappedCommand = append(wrappedCommand, h.Image)
	}

	return append(wrappedCommand, command...)
}

func (h *DockerRunHandler) validate() error {
	if h.Image == "" {
		return errors.New("field image required but not set")
	}
	return nil
}

// InitDockerRunHandler generates a DockerRunHandler
func InitDockerRunHandler(settings map[string]interface{}) (Handler, error) {
	handler := &DockerRunHandler{}
	parsingMetadata := &mapstructure.Metadata{}

	if err := mapstructure.DecodeMetadata(settings, handler, parsingMetadata); err != nil {
		return handler, err
	}
	return handler, validateHandler(handler, parsingMetadata)
}

// ComposeRunHandler wraps a command with `docker-compose exec`
type ComposeRunHandler struct {
	Remove  bool
	Service string
}

// BuildCommand performs the actual wrapping of the command
func (h *ComposeRunHandler) BuildCommand(command []string) []string {
	wrappedCommand := []string{"docker-compose", "run"}

	if h.Remove {
		wrappedCommand = append(wrappedCommand, "--rm")
	}

	if h.Service != "" {
		wrappedCommand = append(wrappedCommand, h.Service)
	}

	return append(wrappedCommand, command...)
}

func (h *ComposeRunHandler) validate() error {
	if h.Service == "" {
		return errors.New("field service required but not set")
	}
	return nil
}

// InitComposeRunHandler generates a ComposeRunHandler
func InitComposeRunHandler(settings map[string]interface{}) (Handler, error) {
	handler := &ComposeRunHandler{}
	parsingMetadata := &mapstructure.Metadata{}

	if err := mapstructure.DecodeMetadata(settings, handler, parsingMetadata); err != nil {
		return handler, err
	}
	return handler, validateHandler(handler, parsingMetadata)
}

// ComposeExecHandler wraps a command with `docker-compose run`
type ComposeExecHandler struct {
	Service string
}

// BuildCommand performs the actual wrapping of the command
func (h *ComposeExecHandler) BuildCommand(command []string) []string {
	wrappedCommand := []string{"docker-compose", "exec"}

	if h.Service != "" {
		wrappedCommand = append(wrappedCommand, h.Service)
	}

	return append(wrappedCommand, command...)
}

func (h *ComposeExecHandler) validate() error {
	if h.Service == "" {
		return errors.New("field service required but not set")
	}
	return nil
}

// InitComposeExecHandler generates a ComposeExecHandler
func InitComposeExecHandler(settings map[string]interface{}) (Handler, error) {
	handler := &ComposeExecHandler{}
	parsingMetadata := &mapstructure.Metadata{}

	if err := mapstructure.DecodeMetadata(settings, handler, parsingMetadata); err != nil {
		return handler, err
	}
	return handler, validateHandler(handler, parsingMetadata)
}

func validateHandler(handler Handler, metadata *mapstructure.Metadata) error {
	if err := handler.validate(); err != nil {
		return err
	}

	if len(metadata.Unused) > 0 {
		return fmt.Errorf("additonal field(s) detected: %v", strings.Join(metadata.Unused, ", "))
	}

	return nil
}
