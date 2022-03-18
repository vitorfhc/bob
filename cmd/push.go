package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushed Docker images",
	Long: `Examples:
  bob push
  bob push --file bobber.yaml`,
	Run: runPush,
}

func runPush(cmd *cobra.Command, args []string) {
	cfg := getBobConfig(cmd)

	logrus.Debug("Running push command with configuration file ", cfg.ConfigPath)

	images, err := cfg.ToImageList()
	if err != nil {
		logrus.WithError(err).Fatal("Error reading configuration file")
	}

	for _, image := range images.Images {
		ctx := context.Background()
		err = image.Push(ctx)
		if err != nil {
			logrus.WithError(err).Fatal("Error pushing image ", image.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
