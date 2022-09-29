package task

import (
	"github.com/kuoss/lethe/file"
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
	var dryRun bool
	deleteBySizeCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "dry run")
	taskCmd.AddCommand(deleteBySizeCmd)
}

func DeleteBySize(cmd *cobra.Command) {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		cmd.PrintErr(err)
		return
	}
	file.DeleteBySize(dryRun)
}
