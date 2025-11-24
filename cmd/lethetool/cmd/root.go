package cmd

import (
	"os"

	"github.com/kuoss/lethe/cmd/lethetool/cmd/version"
	"github.com/kuoss/lethe/cmd/lethetool/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lethetool",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute(version string) {
	// config
	config.Version = version

	// execute
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(version.New())
}
