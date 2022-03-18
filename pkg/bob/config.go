package bob

import (
	"github.com/docker/docker/api/types"
	"github.com/vitorfhc/bob/pkg/docker"
)

// Config struct keeps all configuration variables
type Config struct {
	ConfigPath string
	Debug      bool
	AuthConfig types.AuthConfig
}

// ToImageList reads config.ConfigPath and returns a new ImageList
// generated from the YAML file.
func (c *Config) ToImageList() (*docker.ImageList, error) {
	return docker.NewImageListFromYaml(c.ConfigPath)
}
