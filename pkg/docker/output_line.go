package docker

import "encoding/json"

// OutputLine represents the line of a Docker command output.
type OutputLine struct {
	Error  string `json:"error"`
	Stream string `json:"stream"`
	Status string `json:"status"`
}

func (ol *OutputLine) String() string {
	if ol.Error != "" {
		return ol.Error
	}
	if ol.Stream != "" {
		return ol.Stream
	}
	return ol.Status
}

// HasError returns true if the line has an error.
func (ol *OutputLine) HasError() bool {
	if ol.Error != "" {
		return true
	}
	return false
}

// NewOutputLineFromJSON creates a new ErrorLine from a JSON string
func NewOutputLineFromJSON(jsn string) (*OutputLine, error) {
	outputLine := &OutputLine{}
	err := json.Unmarshal([]byte(jsn), outputLine)
	if err != nil {
		return nil, err
	}
	return outputLine, nil
}
