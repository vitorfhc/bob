package deps

import "github.com/vitorfhc/bob/pkg/docker"

// Manager abstracts a dependencies manager, that knows how to deal with dependencies.
// It could be a graph or simple lists.
type Manager interface {
	AddDependency(from, to *docker.Image) error
	GetDepsFor(image *docker.Image) ([][]*docker.Image, error)
}
