package task

import (
	"github.com/kuoss/lethe/logs"
	"github.com/spf13/cobra"
)

var deleteBySizeCmd = &cobra.Command{
	Use:   "delete-by-size",
	Short: "Delete log files by size",
	Run: func(cmd *cobra.Command, args []string) {
		DeleteBySize(cmd)
	},
}

func initDeleteBySize() {
	taskCmd.AddCommand(deleteBySizeCmd)
}

func DeleteBySize(cmd *cobra.Command) {
	logs.NewRotator().DeleteBySize()
}
