package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitorfhc/bob/pkg/bob"
)

func getBobConfig(cmd *cobra.Command) *bob.Config {
	config := &bob.Config{}

	config.ConfigPath = cmd.Flag("file").Value.String()

	debugString := cmd.Flag("debug").Value.String()
	config.Debug = debugString == "true"

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	usernameFlag := cmd.Flag("username")
	passwordFlag := cmd.Flag("password")
	if usernameFlag != nil && passwordFlag != nil {
		config.AuthConfig.Username = usernameFlag.Value.String()
		config.AuthConfig.Password = passwordFlag.Value.String()
	}

	return config
}
