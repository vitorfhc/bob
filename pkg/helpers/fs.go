package helpers

import (
	"os"
	"path/filepath"
)

// FileExists checks if <filename> is a file and if it exists
func FileExists(filename string) bool {
	stat, err := os.Stat(filename)
	return !os.IsNotExist(err) && !stat.IsDir()
}

// FindFileWithExtensions tries to find if <filename> exists with any of the
// extensions in <extensions> and returns the first one found.
// If no file is found, it returns an empty string and an os.ErrNotExist error.
func FindFileWithExtensions(filename string, extensions []string) (string, error) {
	for _, ext := range extensions {
		if FileExists(filename + ext) {
			return filename + ext, nil
		}
	}
	return "", os.ErrNotExist
}

// GetFileWithoutExt returns the file name without the extension
func GetFileWithoutExt(filename string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
