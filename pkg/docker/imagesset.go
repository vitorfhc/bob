package docker

import "errors"

type ImagesSet struct {
	size    int
	mapping map[string]*Image
}

var globalImagesSet *ImagesSet

func init() {
	globalImagesSet = NewImagesSet()
}

func NewImagesSet(images ...*Image) *ImagesSet {
	set := &ImagesSet{
		mapping: make(map[string]*Image),
	}
	for _, img := range images {
		set.AddImages(img)
	}
	return set
}

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

func (is *ImagesSet) FindImage(id string) (*Image, bool) {
	img, ok := is.mapping[id]
	return img, ok
}

func (is *ImagesSet) IterateImages(f func(string, *Image)) {
	for id, img := range is.mapping {
		f(id, img)
	}
}

func AddImages(imgs ...*Image) error {
	is := globalImagesSet
	return is.AddImages(imgs...)
}

func FindImage(id string) (*Image, bool) {
	is := globalImagesSet
	return is.FindImage(id)
}

func IterateImages(f func(string, *Image)) {
	is := globalImagesSet
	is.IterateImages(f)
}
