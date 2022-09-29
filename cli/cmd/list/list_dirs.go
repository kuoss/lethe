package list

import (
	"bytes"
	"fmt"

	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/kuoss/lethe/file"
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
	dirs := file.ListDirsWithSize()

	var data [][]string
	totalKB := 0
	totalCountFiles := 0
	for _, dir := range dirs {
		totalKB += dir.KB
		totalCountFiles += dir.CountFiles
		firstFile := dir.FirstFile
		lastFile := dir.LastFile
		if firstFile == "" {
			firstFile = "-"
		}
		if lastFile == "" {
			lastFile = "-"
		}
		data = append(data, []string{
			dir.Dirpath,
			fmt.Sprintf("%.1f", float64(dir.KB)/1024),
			fmt.Sprintf("%d", dir.CountFiles),
			firstFile,
			lastFile,
		})
	}
	data = append(data, []string{
		"TOTAL",
		fmt.Sprintf("%.1f", float64(totalKB)/1024),
		fmt.Sprintf("%d", totalCountFiles),
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
