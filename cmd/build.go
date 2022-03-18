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
	Long: `Examples:
  bob build
  bob build --file bobber.yaml
  bob build --file bobber.yaml --file bobber2.yaml`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	logrus.Debug("Running build command with configurations files ", bobPath)

	images, err := docker.NewImageListFromYaml(bobPath)
	if err != nil {
		logrus.WithError(err).Fatal("Error reading configuration file")
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
