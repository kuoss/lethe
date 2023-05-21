package cmd

import (
	"fmt"

	"github.com/kuoss/lethe/letheql"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print LetheQL execution result",
	Run: func(cmd *cobra.Command, args []string) {
		Query(cmd)
	},
}

func init() {
	var query string
	logsCmd.Flags().StringVarP(&query, "query", "q", "", "letheql")
	rootCmd.AddCommand(logsCmd)
}

func Query(cmd *cobra.Command) {
	query, err := cmd.Flags().GetString("query")
	fmt.Println("=== query=", query)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	if query == "" {
		cmd.PrintErr("error: logs command needs an flag: --query\n")
		return
	}
	data, err := letheql.ProcQuery(query, letheql.TimeRange{})
	if err != nil {
		cmd.PrintErr(err)
		return
	}
	for _, log := range data.Logs {
		cmd.Println(log)
	}
}
