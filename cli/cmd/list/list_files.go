package list

import (
	"bytes"
	"fmt"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/logs/rotator"

	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/spf13/cobra"
)

func initListFiles() {
	listCmd.AddCommand(&cobra.Command{
		Use:   "files",
		Short: "List log files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listFiles(cmd)
		},
	})
}

func listFiles(c *cobra.Command) error {
	files, err := rotator.NewRotator().ListFiles()
	if err != nil {
		return fmt.Errorf("error on ListFiles: %w", err)
	}

	var data [][]string
	var totalSize int64
	for _, file := range files {
		totalSize += file.Size
		data = append(data, []string{
			file.FullPath,
			fmt.Sprintf("%.1f", float64(file.Size)/1024/1024),
		})
	}
	data = append(data, []string{"TOTAL", fmt.Sprintf("%.1f", float64(totalSize)/1024/1024)})

	buf := &bytes.Buffer{}
	table := cliutil.NewTableWriter(buf)
	table.SetHeader([]string{"FILEPATH", "SIZE(Mi)"})
	table.AppendBulk(data)
	table.Render()
	c.Print(buf.String())
}
