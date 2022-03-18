package docker

import (
	"github.com/vitorfhc/bob/pkg/helpers/fs"
	"gopkg.in/yaml.v2"
)

// ImageList is a set of Image structs
type ImageList struct {
	Images []Image `yaml:"images"`
}

// NewImageListFromYaml creates a new ImageList from a single YAML config.
func NewImageListFromYaml(file string) (*ImageList, error) {
	content, err := fs.ReadYamlFile(file)
	if err != nil {
		return nil, err
	}

	imageList := &ImageList{}
	err = yaml.Unmarshal(content, imageList)
	if err != nil {
		return nil, err
	}

	return imageList, nil
}
