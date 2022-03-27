package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
	"github.com/vitorfhc/bob/pkg/docker/outputs"
	"github.com/vitorfhc/bob/pkg/helpers/dkr"
)

// ImageConfig has all configuration needed from a YAML file.
type ImageConfig struct {
	Name       string             `yaml:"name"`
	Tags       []string           `yaml:"tags"`
	Context    string             `yaml:"context"`
	Dockerfile string             `yaml:"dockerfile"`
	Target     string             `yaml:"target"`
	BuildArgs  map[string]*string `yaml:"buildArgs"`
	Registry   string             `yaml:"registry"`
}

// Image holds the information about a Docker image.
type Image struct {
	Config *ImageConfig

	buildOnce sync.Once
	logger    *logrus.Entry
}

// NewImage creates an Image using the given configuration.
func NewImage(cfg *ImageConfig) *Image {
	return &Image{
		Config: cfg,
	}
}

// FullName joins the registry with image name
func (i *Image) FullName() string {
	cfg := i.Config
	if cfg.Registry == "" {
		return cfg.Name
	}
	if strings.HasSuffix(cfg.Registry, "/") {
		return cfg.Registry + cfg.Name
	}
	return cfg.Registry + "/" + cfg.Name
}

// Build wraps the internal build function,
// it guarantees to construct the image only once by using sync.Once.
// If the image was already built, it returns (false, nil).
// If any error occurs, it returns (false, err).
// If the image was built succesfully, it returns (true, nil).
func (i *Image) Build() (bool, error) {
	var err error
	built := false
	i.buildOnce.Do(func() {
		ctx := context.Background()
		err = i.build(ctx)
		if err != nil {
			built = true
		}
	})
	return built, err
}

// Build constructs the Docker image.
func (i *Image) build(ctx context.Context) error {
	i.log(logrus.InfoLevel, "Building image", i.FullName())

	cfg := i.Config
	if cfg.Context == "" {
		cfg.Context = "."
	}

	contextPacked, err := archive.TarWithOptions(cfg.Context, &archive.TarOptions{})
	if err != nil {
		return err
	}
	defer contextPacked.Close()

	now := time.Now()
	dockerClient, err := NewClient()
	if err != nil {
		return err
	}
	response, err := dockerClient.ImageBuild(ctx, contextPacked, types.ImageBuildOptions{
		Tags:       i.FullNamesWithTags(),
		Dockerfile: cfg.Dockerfile,
		Target:     cfg.Target,
		BuildArgs:  cfg.BuildArgs,
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

	cfg := i.Config
	dockerClient, err := NewClient()
	if err != nil {
		return err
	}
	for _, tag := range cfg.Tags {
		fullName := i.FullName() + ":" + tag
		i.log(logrus.InfoLevel, "Pushing image", fullName)
		body, err := dockerClient.ImagePush(ctx, fullName, pushOptions)
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
	if len(i.Config.Tags) == 0 {
		return []string{i.FullName() + ":latest"}
	}

	var tags []string
	for _, tag := range i.Config.Tags {
		tags = append(tags, i.FullName()+":"+tag)
	}
	return tags
}

func (i *Image) initLogger() {
	if i.logger == nil {
		i.logger = logrus.WithFields(logrus.Fields{
			"image": i.Config.Name,
		})
	}
}

func (i *Image) log(level logrus.Level, msg ...interface{}) {
	i.initLogger()
	i.logger.Log(level, msg...)
}
