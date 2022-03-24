package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images",
	Long: `Examples:
  bob build
  bob build --file bobber.yaml
  bob build --file bobber.yaml --file bobber2.yaml`,
	Run: runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	cfg := getBobConfig(cmd)

	logrus.Debug("Running build command with configuration file ", cfg.ConfigPath)

	imageList, err := cfg.ToImageList()
	if err != nil {
		logrus.WithError(err).Panic("Error reading configuration file")
	}
	images := imageList.Images

	for _, image := range images {
		ctx := context.Background()
		err = image.Build(ctx)
		if err != nil {
			logrus.WithError(err).Panic("Error building image ", image.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
