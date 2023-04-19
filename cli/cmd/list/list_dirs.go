package list

import (
	"bytes"
	"fmt"

	"github.com/kuoss/lethe/logs/rotator"

	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/spf13/cobra"
)

func initListDirs() {
	listCmd.AddCommand(&cobra.Command{
		Use:   "dirs",
		Short: "List log dirs",
		Run: func(cmd *cobra.Command, args []string) {
			listDirs(cmd)
		},
	})
}

func listDirs(c *cobra.Command) {
	dirs := rotator.NewRotator().ListDirsWithSize()

	var data [][]string
	var totalSize int64
	totalFileCount := 0
	for _, dir := range dirs {
		totalSize += dir.Size
		totalFileCount += dir.FileCount
		firstFile := dir.FirstFile
		lastFile := dir.LastFile
		if firstFile == "" {
			firstFile = "-"
		}
		if lastFile == "" {
			lastFile = "-"
		}
		data = append(data, []string{
			dir.FullPath,
			fmt.Sprintf("%.1f", float64(dir.Size)/1024/1024),
			fmt.Sprintf("%d", dir.FileCount),
			firstFile,
			lastFile,
		})
	}
	data = append(data, []string{
		"TOTAL",
		fmt.Sprintf("%.1f", float64(totalSize)/1024/1024),
		fmt.Sprintf("%d", totalFileCount),
		"-",
		"-",
	})

	buf := &bytes.Buffer{}
	table := cliutil.NewTableWriter(buf)
	table.SetHeader([]string{"DIR", "SIZE(Mi)", "FILES", "FIRST FILE", "LAST FILE"})
	table.AppendBulk(data)
	table.Render()
	c.Print(buf.String())
}
