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
		logrus.WithError(err).Panic("Error reading configuration file")
	}

	for _, image := range images {
		ctx := context.Background()
		err = image.Push(ctx, cfg.AuthConfig)
		if err != nil {
			logrus.WithError(err).Panic("Error pushing image ", image.FullName())
		}
	}
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("username", "u", "", "username for the registry")
	pushCmd.Flags().StringP("password", "p", "", "password for the registry")
	pushCmd.MarkFlagRequired("username")
	pushCmd.MarkFlagRequired("password")
}
