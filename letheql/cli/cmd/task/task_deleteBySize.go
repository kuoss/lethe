package task

import (
	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/logs/rotator"
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
	err := rotator.NewRotator().DeleteBySize()
	if err != nil {
		logger.Errorf("error on DeleteByAge: %s", err)
	}
}
