package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCommands(t *testing.T) {
	cfg, err := GenerateConfig([]byte(fullYaml))
	assert.NoError(t, err)
	assert.ElementsMatch(t, cfg.ListCommands(), []string{"ls", "bundle"})
}

func TestGenerateHandler(t *testing.T) {
	cfg := &Cfg{handler: map[string]Handler{}}
	settings := map[string]interface{}{
		"handler": "docker_compose_run",
		"service": "app",
	}
	err := cfg.generateHandler("test", settings)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.handler["test"])
}

func TestGenerateHandlerError(t *testing.T) {
	cfg := &Cfg{}
	err := cfg.generateHandler("test", map[string]interface{}{
		"handler": "docker_compose_run",
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Error in strategy test: field service required but not set")

	err = cfg.generateHandler("foo", map[string]interface{}{
		"handler": "docker_compose_run",
		"service": "test",
		"other":   "field",
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Error in strategy foo: additonal field(s) detected: other")
}

func TestGetHandlerFor(t *testing.T) {
	var handler Handler
	var err error

	cfg, err := GenerateConfig([]byte(fullYaml))
	assert.NoError(t, err)

	// Without strict mode for a known command
	handler, err = cfg.GetHandlerFor("ls", false)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	// With strict mode for an unkown command
	handler, err = cfg.GetHandlerFor("some-cmd", true)
	assert.Error(t, err, ErrUndefinedCommand)
	assert.Nil(t, handler)

	// Without strict mode for an unkown command
	handler, err = cfg.GetHandlerFor("some-cmd", false)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	// When specifying a path to the command
	handler, err = cfg.GetHandlerFor("/bin/ls", true)
	assert.NoError(t, err)
	assert.NotNil(t, handler)
}
