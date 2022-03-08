package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bob",
	Short: "Bob is an automated tool for building multiple Docker images",
	Long: `Using this tool you may build and push several images in a monorepo.
All you need is a YAML file which has everything you need configured.

Examples:
  bob build
  bob build --file bobber.yaml
  bob build --file bobber.yaml --file bobber2.yaml`,
}

var (
	bobPaths []string
	debug    bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&bobPaths, "file", "f", []string{"bob.yaml"}, "yaml configuration file")
}
