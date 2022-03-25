package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
	"github.com/vitorfhc/bob/pkg/docker/outputs"
	"github.com/vitorfhc/bob/pkg/helpers/dkr"
)

// Image is a struct that represents a Docker image
type Image struct {
	Name       string             `yaml:"name"`
	Tags       []string           `yaml:"tags"`
	Context    string             `yaml:"context"`
	Dockerfile string             `yaml:"dockerfile"`
	Target     string             `yaml:"target"`
	BuildArgs  map[string]*string `yaml:"buildArgs"`
	Registry   string             `yaml:"registry"`

	Built bool

	logger *logrus.Entry
}

// FullName joins the registry with image name
func (i *Image) FullName() string {
	if i.Registry == "" {
		return i.Name
	}
	if strings.HasSuffix(i.Registry, "/") {
		return i.Registry + i.Name
	}
	return i.Registry + "/" + i.Name
}

// Build builds the Docker image
func (i *Image) Build(ctx context.Context) error {
	i.log(logrus.InfoLevel, "Building image", i.FullName())
	if i.Built {
		i.log(logrus.InfoLevel, "Image already built")
		return nil
	}

	if i.Context == "" {
		i.Context = "."
	}

	contextPacked, err := archive.TarWithOptions(i.Context, &archive.TarOptions{})
	if err != nil {
		return err
	}
	defer contextPacked.Close()

	now := time.Now()
	response, err := envClient.ImageBuild(ctx, contextPacked, types.ImageBuildOptions{
		Tags:       i.FullNamesWithTags(),
		Dockerfile: i.Dockerfile,
		Target:     i.Target,
		BuildArgs:  i.BuildArgs,
	})
	if err != nil {
		return err
	}
	defer func() {
		response.Body.Close()
		i.log(logrus.InfoLevel, "Elapsed time", time.Since(now).String())
	}()

	buildOutput := &outputs.BuildOutput{}
	dkr.ScanBody(response.Body, buildOutput, i.logger)

	if buildOutput.HasError() {
		return errors.New(buildOutput.String())
	}

	i.Built = true

	return nil
}

// Push sends the Docker image to the registry
func (i *Image) Push(ctx context.Context, authCfg types.AuthConfig) error {
	authJSON, err := json.Marshal(authCfg)
	if err != nil {
		return err
	}
	authB64 := base64.URLEncoding.EncodeToString(authJSON)
	pushOptions := types.ImagePushOptions{
		RegistryAuth: authB64,
	}

	for _, tag := range i.Tags {
		fullName := i.FullName() + ":" + tag
		i.log(logrus.InfoLevel, "Pushing image", fullName)
		body, err := envClient.ImagePush(ctx, fullName, pushOptions)
		if err != nil {
			return err
		}
		defer body.Close()

		pushOutput := &outputs.PushOutput{}
		dkr.ScanBody(body, pushOutput, i.logger)

		if pushOutput.HasError() {
			return errors.New(pushOutput.String())
		}
	}

	return nil
}

// FullNamesWithTags returns a list of strings, each string is the full name
// of the image with one of its tags.
func (i *Image) FullNamesWithTags() []string {
	if len(i.Tags) == 0 {
		return []string{i.FullName() + ":latest"}
	}

	var tags []string
	for _, tag := range i.Tags {
		tags = append(tags, i.FullName()+":"+tag)
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
	i.logger.Log(level, msg...)
}
