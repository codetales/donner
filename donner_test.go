package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecCommand(t *testing.T) {
	tests := map[string]struct {
		config *Cfg
		params []string
		expErr error
	}{
		"missing command": {
			config: &Cfg{Strategies: map[string]Strategy{}},
			expErr: ErrMissingCommand,
		},
		"undefined command": {
			config: &Cfg{Strategies: map[string]Strategy{"test": {Handler: "test"}}, Commands: map[string]Command{"golang": "exec"}},
			params: []string{"invalid"},
			expErr: ErrUndefinedCommand,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := ExecCommand(test.config, test.params)
			assert.EqualError(t, err, test.expErr.Error())
		})
	}
}
