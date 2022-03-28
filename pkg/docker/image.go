package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vitorfhc/bob/pkg/config"
	"github.com/vitorfhc/bob/pkg/docker/outputs"
	"github.com/vitorfhc/bob/pkg/helpers/dkr"
	"github.com/vitorfhc/bob/pkg/helpers/random"
)

// ImageConfig has all the configurations for a Docker image.
// They are defined in the bob.yaml file.
type ImageConfig struct {
	ID         string             `yaml:"id"`
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
	Config    *ImageConfig
	buildOnce sync.Once
	pushOnce  sync.Once
	logger    *logrus.Entry
}

// NewImage creates an Image using the given configuration.
func NewImage(m map[string]interface{}) (*Image, error) {
	img := &Image{}

	randomID := fmt.Sprintf("%s-%s", random.Adjective(), random.Animal())

	v := viper.New()
	v.SetDefault("id", randomID)
	v.SetDefault("name", "")
	v.SetDefault("tags", []string{"latest"})
	v.SetDefault("context", ".")
	v.SetDefault("dockerfile", "Dockerfile")
	v.SetDefault("target", "")
	v.SetDefault("buildArgs", map[string]*string{})
	v.SetDefault("registry", "")

	for key, val := range m {
		v.Set(key, val)
	}

	fieldsNoEmpty := []string{"id", "name"}
	invalidFields := []string{}
	for _, field := range fieldsNoEmpty {
		value := v.GetString(field)
		if value == "" {
			invalidFields = append(invalidFields, field)
		}
	}
	if len(invalidFields) > 0 {
		errMsg := fmt.Sprintf("Invalid fields: %s", strings.Join(invalidFields, ", "))
		errMsg += ". All these fields are required."
		return nil, errors.New(errMsg)
	}

	err := v.Unmarshal(&img.Config)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// FullName joins the registry with image name
func (i *Image) FullName() string {
	registry := i.Config.Registry
	name := i.Config.Name

	if registry == "" {
		return name
	}

	if strings.HasSuffix(registry, "/") {
		return registry + name
	}

	return registry + "/" + name
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
	i.log(logrus.InfoLevel, "Building image ", i.FullName())

	context := i.Config.Context
	contextPacked, err := archive.TarWithOptions(context, &archive.TarOptions{})
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
		Dockerfile: i.Config.Dockerfile,
		Target:     i.Config.Target,
		BuildArgs:  i.Config.BuildArgs,
	})
	if err != nil {
		return err
	}
	defer func() {
		response.Body.Close()
		i.log(logrus.InfoLevel, "Elapsed time ", time.Since(now).String())
	}()

	buildOutput := &outputs.BuildOutput{}
	dkr.ScanBody(response.Body, buildOutput, i.logger)

	if buildOutput.HasError() {
		return errors.New(buildOutput.String())
	}

	return nil
}

// Push wraps the internal push function, it guarantees to push the image only once by using sync.Once.
// If the image was already pushed, it returns (false, nil).
// If any error occurs, it returns (false, err).
// If the image was pushed succesfully, it returns (true, nil).
func (i *Image) Push() (bool, error) {
	var err error
	var pushed bool
	i.pushOnce.Do(func() {
		ctx := context.Background()
		err = i.push(ctx)
		if err != nil {
			pushed = true
		}
	})
	return pushed, err
}

func (i *Image) push(ctx context.Context) error {
	username := viper.GetString(config.UsernameKey)
	password := viper.GetString(config.PasswordKey)
	auth := types.AuthConfig{
		Username: username,
		Password: password,
	}

	authJSON, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	authB64 := base64.URLEncoding.EncodeToString(authJSON)
	pushOptions := types.ImagePushOptions{
		RegistryAuth: authB64,
	}

	tags := i.Config.Tags
	dockerClient, err := NewClient()
	if err != nil {
		return err
	}
	for _, tag := range tags {
		fullName := i.FullName() + ":" + tag
		i.log(logrus.InfoLevel, "Pushing image ", fullName)
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
	tags := i.Config.Tags
	if len(tags) == 0 {
		return []string{i.FullName() + ":latest"}
	}

	var fullnames []string
	for _, tag := range tags {
		fullnames = append(fullnames, i.FullName()+":"+tag)
	}
	return fullnames
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
