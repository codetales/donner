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
	tests := map[string]struct {
		cfg         *Cfg
		handlerName string
		settings    map[string]interface{}
		expErr      string
	}{
		"valid handler":     {&Cfg{handler: map[string]Handler{}}, "test", map[string]interface{}{"handler": "docker_compose_run", "service": "app",}, ""},
		"invalid handler":   {&Cfg{}, "test", map[string]interface{}{"handler": "docker_compose_run"}, "error in strategy test: field service required but not set"},
		"additional fields": {&Cfg{}, "foo", map[string]interface{}{"handler": "docker_compose_run", "service": "test", "other": "field"}, "error in strategy foo: additonal field(s) detected: other"},
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

func Test_GetHandlerFor(t *testing.T) {
	cfg, err := GenerateConfig([]byte(fullYaml))
	assert.NoError(t, err)

	tests := map[string]struct {
		command    string
		strictMode bool
		expErr     error
	}{
		"known cmd":                    {"ls", false, nil},
		"unknown cmd":                  {"some-cmd", false, nil},
		"unknown cmd, strict":          {"some-cmd", true, ErrUndefinedCommand},
		"custom cmd with path, strict": {"/bin/ls", true, nil},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			handler, err := cfg.GetHandlerFor(test.command, test.strictMode)
			if test.expErr != nil {
				assert.EqualError(t, err, test.expErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, handler)
			}
		})
	}
}
