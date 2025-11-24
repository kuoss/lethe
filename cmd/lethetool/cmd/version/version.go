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
		Long:  "Print the build/version string for lethetool (set via ldflags).",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.Version)
		},
	}
	return cmd
}
