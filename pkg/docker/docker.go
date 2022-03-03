package docker

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Args contains all information needed to create and run
// the docker build command.
type Args struct {
	name       string
	dockerfile string
	tags       []string
	target     string
	buildArgs  map[string]string
}

func (a *Args) generateCommand() *exec.Cmd {
	var args []string

	args = append(args, "--file", a.dockerfile)

	for _, tag := range a.tags {
		imgName := fmt.Sprintf("%s:%s", a.name, tag)
		args = append(args, "--tag", imgName)
	}

	for key, val := range a.buildArgs {
		buildArg := fmt.Sprintf("%s=%s", key, val)
		args = append(args, "--build-arg", buildArg)
	}

	if a.target != "" {
		args = append(args, "--target", a.target)
	}

	cmd := exec.Command("docker", args...)
	logrus.Debug("Generated docker command:", cmd.String())

	return cmd
}

// Run generates the command for the Args and executes it.
func (a *Args) Run() error {
	cmd := a.generateCommand()
	logrus.Info("Building image from", a.dockerfile)
	err := cmd.Run()
	return err
}
