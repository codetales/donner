package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintAliases(t *testing.T) {
	cfg := &Cfg{commands: map[string]Handler{"test": &DockerRunHandler{true, "test", "8080:80", []string{}}}}

	buffer := bytes.Buffer{}
	printAliases(&buffer, cfg, true, true)

	assert.Equal(t, "alias test='donner run --strict --fallback test';\n"+evalInstruction+"#  eval $(donner aliases --strict --fallback)\n", buffer.String())
}
