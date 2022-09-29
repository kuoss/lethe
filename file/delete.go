package file

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
	"github.com/spf13/cast"
)

// DELETE

func DeleteByAge(dryRun bool) {
	w := config.GetWriter()

	retentionTime := config.GetConfig().GetString("retention.time")
	duration, err := util.GetDurationFromAge(retentionTime)
	if err != nil {
		log.Fatalf("cannot parse duration=[%s] error=%s", retentionTime, err)
		return
	}
	point := strings.Replace(config.GetNow().Add(-duration).UTC().String()[0:13], " ", "_", 1)
	files := ListFiles()
	if len(files) < 1 {
		fmt.Fprintf(w, "DeleteByAge( < %s): no files. done.\n", point)
		return
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		if file.Name < point {
			if dryRun {
				fmt.Fprintf(w, "DeleteByAge(%s < %s): %s (dry run)\n", file.Name, point, file.Filepath)
			} else {
				fmt.Fprintf(w, "DeleteByAge(%s < %s): %s\n", file.Name, point, file.Filepath)
				util.Execute("rm -f " + file.Filepath)
			}
		}
	}
	fmt.Fprintf(w, "DeleteByAge(%s): done\n", point)
}

func DeleteBySize(dryRun bool) {
	w := config.GetWriter()

	retentionSizeKB, err := util.StringToKB(config.GetConfig().GetString("retention.size"))
	if err != nil {
		log.Fatal(err)
	}
	files := ListFilesWithSize()
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		diskUsedKB, err := GetDiskUsedKB()
		if err != nil {
			log.Fatalf("cannot find out disk used. error=%s", err)
			return
		}
		if diskUsedKB < retentionSizeKB {
			fmt.Fprintf(w, "DeleteBySize(%.1fm < %.1fm): done\n", float64(diskUsedKB)/1024, float64(retentionSizeKB)/1024)
			return
		}
		if dryRun {
			fmt.Fprintf(w, "DeleteBySize(%.1fm > %.1fm): %s (dry run)\n", float64(diskUsedKB)/1024, float64(retentionSizeKB)/1024, file.Filepath)
		} else {
			fmt.Fprintf(w, "DeleteBySize(%.1fm > %.1fm): %s\n", float64(diskUsedKB)/1024, float64(retentionSizeKB)/1024, file.Filepath)
			util.Execute("rm -f " + file.Filepath)
		}
	}
}

func GetDiskUsedKB() (int, error) {
	kb, _, err := util.Execute("df /data --output=used | tail -1")
	if err != nil {
		return 0, err
	}
	return cast.ToInt(strings.TrimSpace(kb)), nil
}
