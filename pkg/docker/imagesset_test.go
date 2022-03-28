package docker

import "testing"

func TestImagesSet_AddImages(t *testing.T) {
	tests := []struct {
		name        string
		imagesAdded []map[string]interface{}
		imagesToAdd []map[string]interface{}
		wantErr     bool
	}{
		{
			name:        "adds images without error",
			imagesAdded: []map[string]interface{}{},
			imagesToAdd: []map[string]interface{}{
				{"id": "1", "name": "test-image"},
			},
			wantErr: false,
		},
		{
			name: "adds duplicated images",
			imagesAdded: []map[string]interface{}{
				{"id": "1", "name": "test-image"},
			},
			imagesToAdd: []map[string]interface{}{
				{"id": "1", "name": "another-image"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imagesAdded := make([]*Image, len(tt.imagesAdded))
			for i, img := range tt.imagesAdded {
				imagesAdded[i], _ = NewImage(img)
			}
			imagesToAdd := make([]*Image, len(tt.imagesToAdd))
			for i, img := range tt.imagesToAdd {
				imagesToAdd[i], _ = NewImage(img)
			}
			globalImagesSet = NewImagesSet() // reinitializer global ImagesSet
			AddImages(imagesAdded...)
			if err := AddImages(imagesToAdd...); (err != nil) != tt.wantErr {
				t.Errorf("ImagesSet.AddImages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
