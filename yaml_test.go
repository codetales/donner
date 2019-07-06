package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestParseYaml(t *testing.T) {
	yaml, err := parseYaml([]byte(fullYaml))
	expectedResult := &yamlCfg{
		Strategies: map[string]map[string]interface{}{
			"run": {
				"handler": "docker_compose_run",
				"service": "app",
				"remove":  true,
			},
			"run_with_docker": {
				"handler": "docker_run",
				"image":   "alpine:latest",
			},
		},
		DefaultStrategy: "run",
		Commands: map[string]string{
			"ls":     "run_with_docker",
			"bundle": "run",
		},
	}
	assert.NoError(t, err)
	assert.Equal(t, yaml, expectedResult)
}

// TODO: Test parsing yaml with additional fields
