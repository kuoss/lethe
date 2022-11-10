package list

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	cliutil "github.com/kuoss/lethe/cli/util"
	"github.com/kuoss/lethe/file"
	"github.com/spf13/cobra"
)

func initListTargets() {
	listCmd.AddCommand(&cobra.Command{
		Use:   "targets",
		Short: "List targets",
		Run: func(cmd *cobra.Command, args []string) {
			listTargets(cmd)
		},
	})
}

func listTargets(c *cobra.Command) {
	now := time.Now().UTC()
	dirs := file.ListTargets()

	var data [][]string
	var totalSize int64
	totalFileCount := 0
	for _, dir := range dirs {
		totalSize += dir.Size
		totalFileCount += dir.FileCount
		firstFile := dir.FirstFile
		lastFile := dir.LastFile
		lastForward := dir.LastForward
		if firstFile == "" {
			firstFile = "-"
		}
		if lastFile == "" {
			lastFile = "-"
		}
		if lastForward == "" || len(lastForward) != 20 {
			lastForward = "-"
		} else {
			// fmt.Println(lastForward)
			dt, err := time.Parse("2006-01-02T15:04:05Z", lastForward)
			if err != nil {
				lastForward = "-"
			} else {
				age := now.Sub(dt).Round(time.Second)
				re := regexp.MustCompile(`[0-9]+[a-z]`)
				match := re.FindStringSubmatch(age.String())
				lastForward = fmt.Sprintf("%s (%v)", dir.LastForward, match[0])
			}
		}
		data = append(data, []string{
			dir.FullPath,
			fmt.Sprintf("%.1f", float64(dir.Size)/1024/1024),
			fmt.Sprintf("%d", dir.FileCount),
			firstFile,
			lastFile,
			lastForward,
		})
	}
	data = append(data, []string{
		"TOTAL",
		fmt.Sprintf("%.1f", float64(totalSize)/1024/1024),
		fmt.Sprintf("%d", totalFileCount),
		"-",
		"-",
		"-",
	})

	buf := &bytes.Buffer{}
	table := cliutil.NewTableWriter(buf)
	table.SetHeader([]string{"DIR", "SIZE(Mi)", "FILES", "FIRST FILE", "LAST FILE", "LAST FORWARD"})
	table.AppendBulk(data)
	table.Render()
	c.Print(buf.String())
}
