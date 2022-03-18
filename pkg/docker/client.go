package docker

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Client is the Docker client
var Client *client.Client

func init() {
	var err error
	Client, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.WithError(err).Panic("Error creating Docker client")
	}
}
