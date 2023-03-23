package task

import (
	"github.com/kuoss/lethe/logs/rotator"
	"github.com/spf13/cobra"
)

var deleteByAgeCmd = &cobra.Command{
	Use:   "delete-by-age",
	Short: "Delete log files by age",
	Run: func(cmd *cobra.Command, args []string) {
		DeleteByAge(cmd)
	},
}

func initDeleteByAge() {
	taskCmd.AddCommand(deleteByAgeCmd)
}

func DeleteByAge(cmd *cobra.Command) {
	rotator.NewRotator().DeleteByAge()
}
