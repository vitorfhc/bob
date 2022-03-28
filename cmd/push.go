package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitorfhc/bob/pkg/config"
	"github.com/vitorfhc/bob/pkg/docker"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushed Docker images",
	Long: `Examples:
  bob push`,
	Run: runPush,
}

func runPush(cmd *cobra.Command, args []string) {
	config, err := config.ReadConfig()
	if err != nil {
		logrus.WithError(err).Panic("Error reading configuration file")
	}

	logrus.Info("Running push command")

	images := make([]*docker.Image, len(config.Images))
	for i, imageMap := range config.Images {
		image, err := docker.NewImage(imageMap)
		if err != nil {
			logrus.WithError(err).Panic("Error creating image")
		}
		images[i] = image
	}

	for _, image := range images {
		pushed, err := image.Push()
		if err != nil {
			logrus.WithError(err).Panic("Error pushing image")
		}
		if !pushed && err != nil {
			logrus.Info("Image was already pushed successfully")
		} else {
			logrus.Info("Image pushed successfully")
		}
	}
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringP("username", "u", "", "username for the registry")
	viper.BindPFlag(config.UsernameKey, pushCmd.Flags().Lookup("username"))

	pushCmd.Flags().StringP("password", "p", "", "password for the registry")
	viper.BindPFlag(config.PasswordKey, pushCmd.Flags().Lookup("username"))
}
