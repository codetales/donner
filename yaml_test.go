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
    port: "8080:80"
    volumes:
      - "./:/usr/src/app"

default_strategy: run

commands:
  ls: run_with_docker
  bundle: run
`

var yamlWithExtraAttributes = fullYaml + `
something-else:
  foo: bar
`

func TestParseYaml(t *testing.T) {
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
				"port":    "8080:80",
				"volumes": []interface{}{"./:/usr/src/app"},
			},
		},
		DefaultStrategy: "run",
		Commands: map[string]string{
			"ls":     "run_with_docker",
			"bundle": "run",
		},
	}

	yaml, err := parseYaml([]byte(fullYaml))
	assert.NoError(t, err)
	assert.Equal(t, yaml, expectedResult)

	yaml, err = parseYaml([]byte(yamlWithExtraAttributes))
	assert.NoError(t, err)
	assert.Equal(t, yaml, expectedResult)
}
