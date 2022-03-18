package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vitorfhc/bob/pkg/bob"
)

func getBobConfig(cmd *cobra.Command) *bob.Config {
	config := &bob.Config{}

	config.ConfigPath = cmd.Flag("file").Value.String()

	return config
}
