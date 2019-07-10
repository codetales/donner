package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecCommand(t *testing.T) {
	tests := map[string]struct {
		config     *Cfg
		params     []string
		expErr     error
		strictMode bool
	}{
		"missing command": {
			config:     &Cfg{},
			expErr:     ErrMissingCommand,
			strictMode: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := execCommand(test.config, test.params, test.strictMode)
			assert.EqualError(t, err, test.expErr.Error())
		})
	}
}
