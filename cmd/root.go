package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bob",
	Short: "Bob is an automated tool for building multiple Docker images",
	Long: `Using this tool you may build several images in a monorepo.
All you need is a bob.yaml file in the directory which you run the command.
You could also define another file name with -f or --file flag.`,
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {}

func init() {
	rootCmd.PersistentFlags().StringP("file", "f", "bob.yaml", "yaml configuration file")
	rootCmd.PersistentFlags().BoolP("push", "p", false, "push built images")
}
