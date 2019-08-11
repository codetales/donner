package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCommands(t *testing.T) {
	cfg := GenerateConfig()
	err := cfg.Load([]byte(fullYaml))
	assert.NoError(t, err)
	assert.ElementsMatch(t, cfg.ListCommands(), []string{"ls", "bundle"})
}

func TestGenerateHandler(t *testing.T) {
	tests := map[string]struct {
		cfg         *Cfg
		handlerName string
		settings    map[string]interface{}
		expErr      string
	}{
		"valid handler":     {&Cfg{handler: map[string]Handler{}}, "test", map[string]interface{}{"handler": "docker_compose_run", "service": "app"}, ""},
		"invalid handler":   {&Cfg{}, "test", map[string]interface{}{"handler": "docker_compose_run"}, "error in strategy test: field service required but not set"},
		"additional fields": {&Cfg{}, "foo", map[string]interface{}{"handler": "docker_compose_run", "service": "test", "other": "field"}, "error in strategy foo: additional field(s) detected: other"},
		"invalid type":      {&Cfg{}, "foo", map[string]interface{}{"handler": "docker_compose_run", "service": true, "remove": "foo"}, "error in strategy foo: 2 error(s) decoding:\n\n* 'Remove' expected type 'bool', got unconvertible type 'string'\n* 'Service' expected type 'string', got unconvertible type 'bool'"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.cfg.generateHandler(test.handlerName, test.settings)
			if test.expErr != "" {
				assert.EqualError(t, err, test.expErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, test.cfg.handler[test.handlerName])
			}
		})
	}
}

func TestGetHandlerFor(t *testing.T) {
	cfg := GenerateConfig()
	err := cfg.Load([]byte(fullYaml))
	assert.NoError(t, err)

	tests := map[string]struct {
		command      string
		strictMode   bool
		fallbackMode bool
		expErr       error
	}{
		"known cmd":                     {command: "ls"},
		"unknown cmd":                   {command: "some-cmd"},
		"known cmd with path, strict":   {command: "ls", strictMode: true, fallbackMode: true},
		"unknown cmd, strict,":          {command: "some-cmd", strictMode: true, expErr: ErrUndefinedCommand},
		"unknown cmd, strict, fallback": {command: "some-cmd", strictMode: true, expErr: ErrUndefinedCommand},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			handler, err := cfg.GetHandlerFor(test.command, test.strictMode, test.fallbackMode)
			if test.expErr != nil {
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, handler)
			}
		})
	}
}
