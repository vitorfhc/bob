package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunBuildCommand(t *testing.T) {
	args := &Args{
		Name:       "image",
		Tags:       []string{"latest", "1.0.0"},
		Dockerfile: "./testdata/Dockerfile",
		Target:     "",
		BuildArgs:  map[string]string{},
	}
	err := args.Run()
	assert.Nil(t, err, err.Error())
}
