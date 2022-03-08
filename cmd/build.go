package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitorfhc/bob/pkg/helpers"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds Docker images",
	Long:  ``,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	logrus.Debug("Running build command with configuration file: ", bobPath)

	extensions := []string{".yaml", ".yml"}
	withoutExt := helpers.GetFileWithoutExt(bobPath)
	cfgFile, err := helpers.FindFileWithExtensions(withoutExt, extensions)
	if err != nil {
		msg := `Could not find a configuration file with the name %s or any of these extensions %s`
		logrus.WithError(err).Fatalf(msg, bobPath, extensions)
	}
	logrus.Debug("Found configuration file ", cfgFile)

	file, err := os.ReadFile(cfgFile)
	if err != nil {
		logrus.WithError(err).Fatalf("Could not read configuration file %s", cfgFile)
	}

	// images := types.ImageList{}
	// err = yaml.Unmarshal(file, images)
	// if err != nil {
	// 	logrus.WithError(err).Fatalf("Could not parse configuration file %s", cfgFile)
	// }
	// logrus.Info(images)
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
