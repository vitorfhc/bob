package docker

import (
	"bufio"
	"context"
	"errors"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
)

// Image is a struct that represents a Docker image
type Image struct {
	Name       string             `yaml:"name"`
	Tags       []string           `yaml:"tags"`
	Context    string             `yaml:"context"`
	Dockerfile string             `yaml:"dockerfile"`
	Target     string             `yaml:"target"`
	BuildArgs  map[string]*string `yaml:"build_args"`
}

// Build builds the Docker image
func (i *Image) Build(ctx context.Context) error {
	localLogrus := logrus.WithFields(logrus.Fields{
		"image": i.Name,
	})

	localLogrus.Infof("Building image %s", i.Name)

	contextPacked, err := archive.TarWithOptions(i.Context, &archive.TarOptions{})
	if err != nil {
		return err
	}
	defer contextPacked.Close()

	now := time.Now()
	response, err := Client.ImageBuild(ctx, contextPacked, types.ImageBuildOptions{
		Tags:       i.generateFullTags(),
		Dockerfile: i.Dockerfile,
		Target:     i.Target,
		BuildArgs:  i.BuildArgs,
	})
	if err != nil {
		return err
	}
	defer func() {
		response.Body.Close()
		localLogrus.Info("Elapsed time ", time.Since(now))
	}()

	var lastLine string
	var lastLineOutput *OutputLine
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		lastLine = scanner.Text()
		lastLineOutput, _ = NewOutputLineFromJSON(lastLine)
		localLogrus.Debug(lastLineOutput.String())
	}

	if lastLineOutput.HasError() {
		return errors.New(lastLineOutput.String())
	}

	return nil
}

// Push sends the Docker image to the registry
func (i *Image) Push(ctx context.Context) error {
	return nil
}

func (i *Image) generateFullTags() []string {
	if len(i.Tags) == 0 {
		return []string{i.Name}
	}

	var tags []string
	for _, tag := range i.Tags {
		tags = append(tags, i.Name+":"+tag)
	}
	return tags
}
