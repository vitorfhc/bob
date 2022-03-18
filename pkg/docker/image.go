package docker

import (
	"bufio"
	"context"
	"errors"
	"io"
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

	logger *logrus.Entry
}

// Build builds the Docker image
func (i *Image) Build(ctx context.Context) error {
	i.log(logrus.InfoLevel, "Building image ", i.Name)

	contextPacked, err := archive.TarWithOptions(i.Context, &archive.TarOptions{})
	if err != nil {
		return err
	}
	defer contextPacked.Close()

	now := time.Now()
	response, err := envClient.ImageBuild(ctx, contextPacked, types.ImageBuildOptions{
		Tags:       i.generateFullNames(),
		Dockerfile: i.Dockerfile,
		Target:     i.Target,
		BuildArgs:  i.BuildArgs,
	})
	if err != nil {
		return err
	}
	defer func() {
		response.Body.Close()
		i.log(logrus.InfoLevel, "Elapsed time ", time.Since(now).String())
	}()

	lastLineOutput := scanBody(response.Body, i.logger)

	if lastLineOutput.HasError() {
		return errors.New(lastLineOutput.String())
	}

	return nil
}

// Push sends the Docker image to the registry
func (i *Image) Push(ctx context.Context) error {
	for _, tag := range i.Tags {
		fullName := i.Name + ":" + tag
		i.log(logrus.InfoLevel, "Pushing image ", fullName)
		body, err := envClient.ImagePush(ctx, fullName, types.ImagePushOptions{})
		if err != nil {
			return err
		}
		defer body.Close()

		lastLineOutput := scanBody(body, i.logger)

		if lastLineOutput.HasError() {
			return errors.New(lastLineOutput.String())
		}
	}

	return nil
}

func (i *Image) generateFullNames() []string {
	if len(i.Tags) == 0 {
		return []string{i.Name}
	}

	var tags []string
	for _, tag := range i.Tags {
		tags = append(tags, i.Name+":"+tag)
	}
	return tags
}

func (i *Image) initLogger() {
	if i.logger == nil {
		i.logger = logrus.WithFields(logrus.Fields{
			"image": i.Name,
		})
	}
}

func (i *Image) log(level logrus.Level, msg ...interface{}) {
	i.initLogger()
	i.logger.Log(level, msg)
}

func scanBody(body io.ReadCloser, logger *logrus.Entry) *OutputLine {
	var lastLine string
	var lastLineOutput *OutputLine
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		lastLine = scanner.Text()
		lastLineOutput, _ = NewOutputLineFromJSON(lastLine)
		if logger != nil {
			logger.Log(logrus.DebugLevel, lastLineOutput.String())
		}
	}
	return lastLineOutput
}
