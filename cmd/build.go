package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitorfhc/bob/pkg/config"
	"github.com/vitorfhc/bob/pkg/docker"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images",
	Long: `Examples:
  bob build`,
	Run: runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	config, err := config.ReadConfig()
	if err != nil {
		logrus.WithError(err).Panic("Error reading configuration file")
	}

	logrus.Info("Running build command")

	images := make([]*docker.Image, len(config.Images))
	for i, imageMap := range config.Images {
		image, err := docker.NewImage(imageMap)
		if err != nil {
			logrus.WithError(err).Panic("Error creating image")
		}
		images[i] = image
	}

	for _, image := range images {
		built, err := image.Build()
		if err != nil {
			logrus.WithError(err).Panic("Error building image")
		}
		if !built && err != nil {
			logrus.Info("Image was already built successfully")
		} else {
			logrus.Info("Image built successfully")
		}
	}
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
