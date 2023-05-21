package list

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List log files or dirs",
}

func Init(parentCmd *cobra.Command) {
	initListDirs()
	initListFiles()
	initListTargets()
	parentCmd.AddCommand(listCmd)
}

// for test
func GetListCmd() *cobra.Command {
	return listCmd
}
