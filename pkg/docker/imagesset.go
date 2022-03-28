package docker

import "errors"

// ImagesSet is a set of images. It maps image IDs to images.
type ImagesSet struct {
	size    int
	mapping map[string]*Image
}

var globalImagesSet *ImagesSet

func init() {
	globalImagesSet = NewImagesSet()
}

// NewImagesSet returns a new instance of ImagesSet.
// You may pass initial images to the constructor.
func NewImagesSet(images ...*Image) *ImagesSet {
	set := &ImagesSet{
		mapping: make(map[string]*Image),
	}
	for _, img := range images {
		set.AddImages(img)
	}
	return set
}

// AddImages adds images to the set.
// It returns an error if there are any ID collisions.
func (is *ImagesSet) AddImages(imgs ...*Image) error {
	for _, img := range imgs {
		_, ok := is.mapping[img.Config.ID]
		if ok {
			return errors.New("image with this ID already exists in ImagesSet")
		}
		is.mapping[img.Config.ID] = img
		is.size++
	}
	return nil
}

// FindImage searches for an image ID and returns (*Image, true)
// if found. If not found it returns (nil, false).
func (is *ImagesSet) FindImage(id string) (*Image, bool) {
	img, ok := is.mapping[id]
	return img, ok
}

// IterateImages iterates over all images in the set.
func (is *ImagesSet) IterateImages(f func(string, *Image)) {
	for id, img := range is.mapping {
		f(id, img)
	}
}

// AddImages adds images to the global set.
// It returns an error if there are any ID collisions.
func AddImages(imgs ...*Image) error {
	is := globalImagesSet
	return is.AddImages(imgs...)
}

// FindImage searches in the global set for an image ID and returns (*Image, true)
// if found. If not found it returns (nil, false).
func FindImage(id string) (*Image, bool) {
	is := globalImagesSet
	return is.FindImage(id)
}

// IterateImages iterates over all images in the global set.
func IterateImages(f func(string, *Image)) {
	is := globalImagesSet
	is.IterateImages(f)
}
