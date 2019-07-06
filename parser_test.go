package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestGenerationErrors {
// 	tests := map[string]struct {
// 		input  string
// 		exp    *Cfg
// 		expErr error
// 	}{
// 		"yaml without commands":   {input: missingCommandsYaml, expErr: ErrNoCommandsSpecified},
// 		"yaml without strategies": {input: missingStrategiesYml, expErr: ErrNoStrategiesSpecified},
// 	}

// 	for name, test := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			_, err := parseYaml([]byte(test.input))
// 			assert.Error(t, err, test.expErr.Error())
// 		})
// 	}
// }

func TestListCommands(t *testing.T) {
	cfg, err := generateConfig([]byte(fullYaml))
	assert.NoError(t, err)
	assert.ElementsMatch(t, cfg.ListCommands(), []string{"ls", "bundle"})
}

func TestGetHandlerFor(t *testing.T) {
	var handler CommandWrapper
	var err error

	cfg, err := generateConfig([]byte(fullYaml))
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
