package docker

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Args contains all information needed to create and run
// the docker build command.
type Args struct {
	Name       string
	Tags       []string
	Dockerfile string
	Target     string
	BuildArgs  map[string]string
	Context    string
}

func (a *Args) generateCommand() *exec.Cmd {
	var args []string

	args = append(args, "build")

	args = append(args, "--file", a.Dockerfile)

	for _, tag := range a.Tags {
		imgName := fmt.Sprintf("%s:%s", a.Name, tag)
		args = append(args, "--tag", imgName)
	}

	for key, val := range a.BuildArgs {
		buildArg := fmt.Sprintf("%s=%s", key, val)
		args = append(args, "--build-arg", buildArg)
	}

	if a.Target != "" {
		args = append(args, "--target", a.Target)
	}

	args = append(args, a.Context)

	cmd := exec.Command("docker", args...)
	logrus.Info("Generated docker command: ", cmd.String())

	return cmd
}

// Run generates the command for the Args and executes it.
func (a *Args) Run() error {
	logrus.Info("Building image from ", a.Dockerfile)
	cmd := a.generateCommand()
	err := cmd.Run()
	return err
}
