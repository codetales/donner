package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDockerRunHandler_BuildCommand(t *testing.T) {
	want := []string{"docker", "run", "-it", "--rm", "test", "-p ", "8080:80", "-v ", "./:/opt/", "ls -la"}
	dockerRunHandler := &DockerRunHandler{true, "test", "8080:80", []string{"./:/opt/"}}

	got := dockerRunHandler.BuildCommand([]string{"ls -la"})

	assert.Equal(t, want, got)
}

func TestComposeExecHandler_BuildCommand(t *testing.T) {
	want := []string{"docker-compose", "exec", "gotest", "ls -la"}
	composeExecHandler := &ComposeExecHandler{"gotest"}

	got := composeExecHandler.BuildCommand([]string{"ls -la"})

	assert.Equal(t, want, got)
}

func TestComposeRunHandler_BuildCommand(t *testing.T) {
	want := []string{"docker-compose", "run", "--rm", "gotest", "ls -la"}
	composeRunHandler := &ComposeRunHandler{true, "gotest"}

	got := composeRunHandler.BuildCommand([]string{"ls -la"})

	assert.Equal(t, want, got)
}
