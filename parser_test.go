package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = `
strategies:
  run:
    handler: docker_compose_run
    service: app
    remove: true
  exec:
    handler: docker_compose_exec
    service: app
`

func TestParsing(t *testing.T) {
	// TODO add many more tests
	tests := map[string]struct {
		input  string
		exp    *Cfg
		expErr error
	}{
		"partial valid yaml": {data, &Cfg{Strategies: map[string]Strategy{"exec": Strategy{Handler: "docker_compose_exec", Service: "app", Remove: false}, "run": Strategy{Handler: "docker_compose_run", Service: "app", Remove: true}}, DefaultStrategy: "", Commands: map[string]Command(nil)}, nil},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := ParseFile([]byte(test.input))
			if test.expErr != nil {
				assert.Error(t, err, test.expErr.Error())
			} else {
				assert.Equal(t, test.exp, res)
			}
		})
	}
}
