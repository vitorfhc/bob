package docker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCommand(t *testing.T) {
	args := &Args{
		name:       "nginx",
		tags:       []string{"latest", "1.0.0"},
		dockerfile: "./Dockerfile",
		target:     "",
		buildArgs:  map[string]string{},
	}

	cmd := args.generateCommand()
	dockerPath := cmd.Path

	expected := fmt.Sprintf("%s --file ./Dockerfile --tag nginx:latest --tag nginx:1.0.0", dockerPath)
	assert.Equal(t, expected, cmd.String())

	// Adds build args
	args.buildArgs = map[string]string{
		"ARG1": "VAL1",
		"ARG2": "VAL2",
	}
	cmd = args.generateCommand()
	expected = expected + " --build-arg ARG1=VAL1"
	expected = expected + " --build-arg ARG2=VAL2"
	assert.Equal(t, expected, cmd.String())

	// Adds target
	args.target = "final"
	cmd = args.generateCommand()
	expected = expected + " --target final"
	assert.Equal(t, expected, cmd.String())
}
