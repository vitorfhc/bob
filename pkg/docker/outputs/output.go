package outputs

// Output represents the line of a Docker command output.
type Output interface {
	String() string
	HasError() bool
	LoadFromJSON(string) error
}
