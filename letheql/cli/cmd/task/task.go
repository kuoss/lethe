package task

import (
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Run task",
}

func Init(parentCmd *cobra.Command) {
	initDeleteByAge()
	initDeleteBySize()
	parentCmd.AddCommand(taskCmd)
}
