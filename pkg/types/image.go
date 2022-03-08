package types

// Image is a struct that represents a Docker image
type Image struct {
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags"`
}

// ImageList is a set of Image structs
type ImageList struct {
	Images []Image `yaml:"images"`
}
