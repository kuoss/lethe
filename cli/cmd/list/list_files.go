package list

import (
	"bytes"
	"fmt"

	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/kuoss/lethe/file"
	"github.com/spf13/cobra"
)

func initListFiles() {
	listCmd.AddCommand(&cobra.Command{
		Use:   "files",
		Short: "List log files",
		Run: func(cmd *cobra.Command, args []string) {
			listFiles(cmd)
		},
	})
}

func listFiles(c *cobra.Command) {
	files := file.ListFilesWithSize()

	var data [][]string
	total := 0
	for _, file := range files {
		total += file.KB
		data = append(data, []string{
			file.Filepath,
			fmt.Sprintf("%.1f", float64(file.KB)/1024),
		})
	}
	data = append(data, []string{"TOTAL", fmt.Sprintf("%.1f", float64(total)/1024)})

	buf := &bytes.Buffer{}
	table := cliutil.NewTableWriter(buf)
	table.SetHeader([]string{"FILEPATH", "SIZE(Mi)"})
	table.AppendBulk(data)
	table.Render()
	c.Print(buf.String())
}
