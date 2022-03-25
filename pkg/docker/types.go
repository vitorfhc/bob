package docker

import (
	"sync"

	"github.com/sirupsen/logrus"
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
