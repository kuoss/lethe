package version

import (
	"fmt"

	"github.com/kuoss/lethe/cmd/lethetool/config"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the lethetool version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.Version)
		},
	}
	return cmd
}
