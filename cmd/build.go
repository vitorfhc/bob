package cmd

import (
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

	for _, image := range imageList {
		built, err := image.Build()
		if err != nil {
			logrus.WithError(err).Panic("Error building image ", image.FullName())
		}
		if !built {
			logrus.Infof("Image %s was already built", image.FullName())
		}
	}
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
