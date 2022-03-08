package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitorfhc/bob/pkg/docker"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images",
	Long:  ``,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	logrus.Debug("Running build command with configurations files ", bobPaths)

	images, err := docker.NewImageListFromYamls(bobPaths...)
	if err != nil {
		logrus.WithError(err).Fatal("Error reading configuration files")
	}

	for _, image := range images.Images {
		ctx := context.Background()
		err = image.Build(ctx)
		if err != nil {
			logrus.WithError(err).Fatal("Error building image ", image.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
