package outputs

import "encoding/json"

// PushOutput represents the line of a Docker push command output.
type PushOutput struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (bp *PushOutput) String() string {
	if bp.Error != "" {
		return bp.Error
	}
	return bp.Status
}

// HasError returns true if PushOutput has an error.
func (bp *PushOutput) HasError() bool {
	if bp.Error != "" {
		return true
	}
	return false
}

// LoadFromJSON loads the PushOutput from a JSON string.
func (bp *PushOutput) LoadFromJSON(j string) error {
	err := json.Unmarshal([]byte(j), bp)
	if err != nil {
		return err
	}
	return nil
}
