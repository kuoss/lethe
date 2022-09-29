package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of lethetool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("lethetool v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
