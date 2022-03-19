package fs

import (
	"io/ioutil"
	"os"
	"strings"
)

// FileExists checks if <filename> is a file and if it exists.
func FileExists(filename string) bool {
	stat, err := os.Stat(filename)
	return !os.IsNotExist(err) && !stat.IsDir()
}

// FindFileWithExtensions tries to find if <filename> exists with any of the
// extensions in <extensions> and returns the first one found.
// If no file is found, it returns an empty string and an os.ErrNotExist error.
func FindFileWithExtensions(filename string, extensions []string) (string, error) {
	for _, ext := range extensions {
		searchFilename := filename
		if strings.HasSuffix(filename, ext) {
			searchFilename = filename[:len(filename)-len(ext)]
		}
		if FileExists(searchFilename + ext) {
			return searchFilename + ext, nil
		}
	}
	return "", os.ErrNotExist
}

// ReadYamlFile tries [.yml, .yaml] extensions to read a file and returns its
// contents as a byte slice.
func ReadYamlFile(filename string) ([]byte, error) {
	file, err := FindFileWithExtensions(filename, []string{".yaml", ".yml"})
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(file)
}
