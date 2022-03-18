package docker

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

var envClient *client.Client

func init() {
	var err error
	envClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.WithError(err).Panic("Error creating Docker client")
	}
}
