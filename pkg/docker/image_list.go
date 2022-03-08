package docker

import (
	"github.com/vitorfhc/bob/pkg/helpers"
	"gopkg.in/yaml.v2"
)

// ImageList is a set of Image structs
type ImageList struct {
	Images []Image `yaml:"images"`
}

// NewImageListFromYamls creates a new ImageList from YAML configs.
// The last YAML config has the highest priority and will override the rest.
func NewImageListFromYamls(files ...string) (*ImageList, error) {
	imageLists := make([]ImageList, len(files))

	for ind, file := range files {
		content, err := helpers.ReadYamlFile(file)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(content, &imageLists[ind])
		if err != nil {
			return nil, err
		}
	}

	imageList := &ImageList{}

	for _, list := range imageLists {
		for _, image := range list.Images {
			imageList.Images = append(imageList.Images, image)
		}
	}

	return imageList, nil
}
