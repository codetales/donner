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

}
