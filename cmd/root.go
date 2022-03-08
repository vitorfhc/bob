package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bob",
	Short: "Bob is an automated tool for building multiple Docker images",
	Long: `Using this tool you may build several images in a monorepo.
All you need is a bob.yaml or bob.yml file in the directory which you run the command.
You could also define another file name with -f or --file flag.`,
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
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	logrus.SetLevel(logrus.DebugLevel)
}
