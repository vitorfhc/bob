package docker

import (
	"reflect"
	"testing"
)

func TestImage_generateFullNames(t *testing.T) {
	type fields struct {
		Name string
		Tags []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "image with no tags",
			fields: fields{
				Name: "test-image",
				Tags: []string{},
			},
			want: []string{"test-image:latest"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Image{
				Name: tt.fields.Name,
				Tags: tt.fields.Tags,
			}
			if got := i.generateFullNames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.generateFullNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
