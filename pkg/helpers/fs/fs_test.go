package fs

import (
	"testing"
)

func TestFindFileWithExtensions(t *testing.T) {
	type args struct {
		filename   string
		extensions []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "inexistent file",
			args: args{
				filename:   "testdada/inexistent",
				extensions: []string{".yaml", ".yml"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "existent file with .yml extension",
			args: args{
				filename:   "testdata/test_01",
				extensions: []string{".yaml", ".yml"},
			},
			want:    "testdata/test_01.yml",
			wantErr: false,
		},
		{
			name: "existent file with .yaml extension",
			args: args{
				filename:   "testdata/test_02",
				extensions: []string{".yaml", ".yml"},
			},
			want:    "testdata/test_02.yaml",
			wantErr: false,
		},
		{
			name: "existent file all extensions",
			args: args{
				filename:   "testdata/test_03",
				extensions: []string{".yaml", ".yml"},
			},
			want:    "testdata/test_03.yaml",
			wantErr: false,
		},
		{
			name: "existent file but input with .yml extension",
			args: args{
				filename:   "testdata/test_01.yml",
				extensions: []string{".yaml", ".yml"},
			},
			want:    "testdata/test_01.yml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindFileWithExtensions(tt.args.filename, tt.args.extensions)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFileWithExtensions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindFileWithExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "existent file",
			args: args{
				filename: "testdata/test_01.yml",
			},
			want: true,
		},
		{
			name: "inexistent file",
			args: args{
				filename: "testdata/test_02.yml",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.filename); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
