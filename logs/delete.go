package logs

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
)

// DELETE
func (rotator *Rotator) DeleteByAge() {
	retentionTime := config.GetConfig().GetString("retention.time")
	duration, err := util.GetDurationFromAge(retentionTime)
	if err != nil {
		log.Fatalf("cannot parse duration=[%s] error=%s", retentionTime, err)
		return
	}
	point := strings.Replace(clock.GetNow().Add(-duration).UTC().String()[0:13], " ", "_", 1)
	files := rotator.ListFiles()
	if len(files) < 1 {
		fmt.Printf("DeleteByAge( < %s): no files. done.\n", point)
		return
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		if file.Name < point {
			fmt.Printf("DeleteByAge(%s < %s): %s\n", file.Name, point, file.FullPath)
			err := rotator.driver.Delete(file.FullPath)
			if err != nil {
				return
			}
		}
	}
	fmt.Printf("DeleteByAge(%s): Done\n", point)
}

func (rotator *Rotator) DeleteBySize() {
	retentionSizeBytes, err := util.StringToBytes(config.GetConfig().GetString("retention.size"))
	if err != nil {
		log.Fatal(err)
	}
	files := rotator.ListFilesWithSize()
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		usedBytes, err := rotator.GetUsedBytes(config.GetLogRoot())
		if err != nil {
			log.Fatalf("cannot get used bytes: %s", err)
			return
		}
		if usedBytes < retentionSizeBytes {
			fmt.Printf("DeleteBySize(%d < %d): Done\n", usedBytes, retentionSizeBytes)
			return
		}
		fmt.Printf("DeleteBySize(%d > %d): %s\n", usedBytes, retentionSizeBytes, file.FullPath)
		rotator.driver.Delete(file.FullPath)
	}
}
