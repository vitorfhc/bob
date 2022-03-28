package docker

import (
	"reflect"
	"testing"
)

func TestImage_FullNamesWithTags(t *testing.T) {
	type fields map[string]interface{}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "image with no tags",
			fields: map[string]interface{}{
				"name": "test-image",
				"tags": []string{},
			},
			want: []string{"test-image:latest"},
		},
		{
			name: "image with one tag",
			fields: map[string]interface{}{
				"name": "test-image",
				"tags": []string{"test-tag"},
			},
			want: []string{"test-image:test-tag"},
		},
		{
			name: "image with multiple tags",
			fields: map[string]interface{}{
				"name": "test-image",
				"tags": []string{"test-tag", "test-tag-2"},
			},
			want: []string{"test-image:test-tag", "test-image:test-tag-2"},
		},
	}
	for _, tt := range tests {
		img, err := NewImage(tt.fields)
		t.Run(tt.name, func(t *testing.T) {
			if (err != nil) != tt.wantErr {
				t.Errorf("NewImage() error = %v, wantErr %v", err, tt.wantErr)
				t.FailNow()
			}
			if got := img.FullNamesWithTags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Image.FullNamesWithTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_FullName(t *testing.T) {
	type fields map[string]interface{}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "image with double slash registry",
			fields: map[string]interface{}{
				"name":     "test-image",
				"registry": "test-registry/",
			},
			want: "test-registry/test-image",
		},
		{
			name: "image with no registry",
			fields: map[string]interface{}{
				"name": "test-image",
			},
			want: "test-image",
		},
		{
			name: "image with registry",
			fields: map[string]interface{}{
				"name":     "test-image",
				"registry": "test-registry",
			},
			want: "test-registry/test-image",
		},
	}
	for _, tt := range tests {
		img, err := NewImage(tt.fields)
		t.Run(tt.name, func(t *testing.T) {
			if (err != nil) != tt.wantErr {
				t.Errorf("NewImage() error = %v, wantErr %v", err, tt.wantErr)
				t.FailNow()
			}
			if got := img.FullName(); got != tt.want {
				t.Errorf("Image.FullName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewImage(t *testing.T) {
	type fields map[string]interface{}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "image with defaults",
			fields: map[string]interface{}{
				"id":   "test-id", // avoids random id
				"name": "test-image",
			},
			want: map[string]interface{}{
				"id":         "test-id",
				"name":       "test-image",
				"tags":       []string{"latest"},
				"context":    ".",
				"dockerfile": "Dockerfile",
				"target":     "",
				"buildArgs":  map[string]*string{},
				"registry":   "",
			},
		},
		{
			name: "image with no defaults",
			fields: map[string]interface{}{
				"id":         "test-id",
				"name":       "xablau",
				"tags":       []string{"lol"},
				"context":    "./dir",
				"dockerfile": "duckerfil",
				"target":     "there",
				"buildArgs": map[string]*string{
					"foo": nil,
				},
				"registry": "duckuhub",
			},
			want: map[string]interface{}{
				"id":         "test-id",
				"name":       "xablau",
				"tags":       []string{"lol"},
				"context":    "./dir",
				"dockerfile": "duckerfil",
				"target":     "there",
				"buildArgs": map[string]*string{
					"foo": nil,
				},
				"registry": "duckuhub",
			},
		},
		{
			name: "image with no name",
			fields: map[string]interface{}{
				"id": "test-id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := NewImage(tt.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewImage() error = %v, wantErr %v", err, tt.wantErr)
				t.FailNow()
			}
			if err != nil {
				return
			}
			got := map[string]interface{}{
				"id":         img.Config.ID,
				"name":       img.Config.Name,
				"tags":       img.Config.Tags,
				"context":    img.Config.Context,
				"dockerfile": img.Config.Dockerfile,
				"target":     img.Config.Target,
				"buildArgs":  img.Config.BuildArgs,
				"registry":   img.Config.Registry,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
