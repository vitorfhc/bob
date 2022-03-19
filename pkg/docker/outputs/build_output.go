package outputs

import "encoding/json"

// BuildOutput represents the line of a Docker build command output.
type BuildOutput struct {
	Stream string `json:"stream"`
	Error  string `json:"error"`
}

func (bo *BuildOutput) String() string {
	if bo.Error != "" {
		return bo.Error
	}
	return bo.Stream
}

// HasError returns true if BuildOutput has an error.
func (bo *BuildOutput) HasError() bool {
	if bo.Error != "" {
		return true
	}
	return false
}

// LoadFromJSON loads the BuildOutput from a JSON string.
func (bo *BuildOutput) LoadFromJSON(j string) error {
	err := json.Unmarshal([]byte(j), bo)
	if err != nil {
		return err
	}
	return nil
}
