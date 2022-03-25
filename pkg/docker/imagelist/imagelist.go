package imagelist

import (
	"github.com/samber/lo"
	"github.com/vitorfhc/bob/pkg/docker"
	"github.com/vitorfhc/bob/pkg/helpers/fs"
	"gopkg.in/yaml.v2"
)

type imagesConfigsList struct {
	ImagesConfigs []docker.ImageConfig `yaml:"images"`
}

// ImageList is an alias for a slice of *docker.Image
type ImageList []*docker.Image

// NewFromYaml creates a new ImageList from a single YAML config.
func NewFromYaml(file string) (ImageList, error) {
	content, err := fs.ReadYamlFile(file)
	if err != nil {
		return nil, err
	}

	icl := &imagesConfigsList{}
	err = yaml.Unmarshal(content, icl)
	if err != nil {
		return nil, err
	}

	images := make(ImageList, len(icl.ImagesConfigs))
	lo.ForEach(icl.ImagesConfigs, func(config docker.ImageConfig, _ int) {
		img := docker.NewImage(&config)
		images = append(images, img)
	})

	return images, nil
}
