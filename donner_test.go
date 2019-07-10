package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecCommand(t *testing.T) {
	tests := map[string]struct {
		config       *Cfg
		params       []string
		expErr       error
		strictMode   bool
		fallbackMode bool
	}{
		"missing command": {
			config:       &Cfg{},
			expErr:       ErrMissingCommand,
			strictMode:   true,
			fallbackMode: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := execCommand(test.config, test.params, test.strictMode, test.fallbackMode)
			assert.EqualError(t, err, test.expErr.Error())
		})
	}
}
