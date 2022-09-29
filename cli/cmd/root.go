package cmd

import (
	"github.com/kuoss/lethe/cli/cmd/list"
	"github.com/kuoss/lethe/cli/cmd/task"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "lethetool",
		Short: "Tooling for the Lethe logging system.",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	list.Init(rootCmd)
	task.Init(rootCmd)
}

// for test
func GetRootCmd() *cobra.Command {
	return rootCmd
}
