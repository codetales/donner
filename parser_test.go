package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var missingCommandsYaml = `
strategies:
  exec:
    handler: docker_compose_exec
    service: app
`

var missingStrategiesYml = `
commands:
  rails: exec
  rspec: exec
`

var fullYaml = `
strategies:
  run:
    handler: docker_compose_run
    service: app
    remove: true
  run_with_docker:
    handler: docker_run
    image: alpine:latest

default_strategy: run

commands:
  ls: run_with_docker
  bundle: run
`

func TestParseFile(t *testing.T) {
	tests := map[string]struct {
		input  string
		exp    *Cfg
		expErr error
	}{
		"yaml without commands":   {input: missingCommandsYaml, expErr: ErrNoCommandsSpecified},
		"yaml without strategies": {input: missingStrategiesYml, expErr: ErrNoStrategiesSpecified},
		"full yaml spec":          {input: fullYaml, exp: &Cfg{Strategies: map[string]Strategy{"run": {Handler: "docker_compose_run", Service: "app", Remove: true}, "run_with_docker": {Handler: "docker_run", Image: "alpine:latest"}}, DefaultStrategy: "run", Commands: map[string]Command{"ls": "run_with_docker", "bundle": "run"}}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := parseFile([]byte(test.input))
			if test.expErr != nil {
				assert.Error(t, err, test.expErr.Error())
			} else {
				assert.Equal(t, test.exp, res)
			}
		})
	}
}

func TestListCommands(t *testing.T) {
	cfg, err := parseFile([]byte(fullYaml))
	assert.NoError(t, err)
	assert.ElementsMatch(t, cfg.ListCommands(), []string{"ls", "bundle"})
}
