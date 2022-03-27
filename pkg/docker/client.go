package docker

import (
	"github.com/docker/docker/client"
)

// NewClient wraps the internal client.NewClientWithOpts function,
// which creates a new Docker client using the options provided.
// The default options used are the ones provided by client.FromEnv,
// and the API version negotiation is enabled.
func NewClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
