package bob

import (
	"github.com/docker/docker/api/types"
	"github.com/vitorfhc/bob/pkg/docker/imagelist"
)

// Config struct keeps all configuration variables
type Config struct {
	ConfigPath string
	Debug      bool
	AuthConfig types.AuthConfig
}

// ToImageList reads config.ConfigPath and returns a new ImageList
// generated from the YAML file.
func (c *Config) ToImageList() (imagelist.ImageList, error) {
	return imagelist.NewFromYaml(c.ConfigPath)
}
