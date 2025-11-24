package cmd

import (
	"os"

	"github.com/kuoss/lethe/cmd/lethetool/cmd/version"
	"github.com/kuoss/lethe/cmd/lethetool/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lethetool",
	Short: "Tooling for the Lethe logging system.",
}

func Execute(ver string) {
	config.Version = ver

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(version.New())
}
