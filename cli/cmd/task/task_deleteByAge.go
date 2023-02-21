package task

import (
	"github.com/kuoss/lethe/logs"
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
	var dryRun bool
	deleteByAgeCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "dry run")
	taskCmd.AddCommand(deleteByAgeCmd)
}

func DeleteByAge(cmd *cobra.Command) {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		cmd.PrintErr(err)
		return
	}
	logs.DeleteByAge(dryRun)
}
