package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vitorfhc/bob/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "bob",
	Short: "Bob is an automated tool for building multiple Docker images",
	Long: `Using this tool you may build and push several images in a monorepo.
All you need is a YAML file which has everything you need configured.

Examples:
  bob build
  bob push`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "wether to print debug messages")
	viper.BindPFlag(config.DebugKey, rootCmd.PersistentFlags().Lookup("debug"))
}
