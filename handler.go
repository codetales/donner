package main

import "github.com/mitchellh/mapstructure"

type DockerRunHandler struct {
	Remove bool
	Image  string
}

func (handler *DockerRunHandler) Validate() error {
	return nil
}

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

func InitDockerRunHandler(settings map[string]interface{}) (Handler, error) {
	var handler *DockerRunHandler
	if err := mapstructure.Decode(settings, &handler); err != nil {
		return handler, err
	}

	if err := handler.Validate(); err != nil {
		return handler, err
	}

	return handler, nil
}

type ComposeRunHandler struct {
	Remove  bool
	Service string
}

func (handler *ComposeRunHandler) Validate() error {
	return nil
}

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

func InitComposeRunHandler(settings map[string]interface{}) (Handler, error) {
	var handler *ComposeRunHandler
	if err := mapstructure.Decode(settings, &handler); err != nil {
		return handler, err
	}

	if err := handler.Validate(); err != nil {
		return handler, err
	}

	return handler, nil
}

type ComposeExecHandler struct {
	Service string
}

func (handler *ComposeExecHandler) Validate() error {
	return nil
}

func (handler *ComposeExecHandler) WrapCommand(command []string) []string {
	wrappedCommand := []string{"docker-compose", "run"}

	if handler.Service != "" {
		wrappedCommand = append(wrappedCommand, handler.Service)
	}

	return append(wrappedCommand, command...)
}

func InitComposeExecHandler(settings map[string]interface{}) (Handler, error) {
	var handler *ComposeExecHandler
	if err := mapstructure.Decode(settings, &handler); err != nil {
		return handler, err
	}

	if err := handler.Validate(); err != nil {
		return handler, err
	}

	return handler, nil
}
