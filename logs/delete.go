package logs

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
)

// DELETE
func (rotator *Rotator) DeleteByAge(dryRun bool) {
	w := config.GetWriter()

	retentionTime := config.GetConfig().GetString("retention.time")
	duration, err := util.GetDurationFromAge(retentionTime)
	if err != nil {
		log.Fatalf("cannot parse duration=[%s] error=%s", retentionTime, err)
		return
	}
	point := strings.Replace(config.GetNow().Add(-duration).UTC().String()[0:13], " ", "_", 1)
	files := rotator.ListFiles()
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
				//		fmt.Fprintf(w, "DeleteByAge(%s < %s): %s (dry run)\n", file.Name, point, file.FullPath)
			} else {
				fmt.Fprintf(w, "DeleteByAge(%s < %s): %s\n", file.Name, point, file.FullPath)
				err := rotator.driver.Delete(file.FullPath)
				if err != nil {
					return
				}
			}
		}
	}
	fmt.Fprintf(w, "DeleteByAge(%s): done\n", point)
}

func (rotator *Rotator) DeleteBySize(dryRun bool) {
	//	w := config.GetWriter()

	retentionSizeKB, err := util.StringToKB(config.GetConfig().GetString("retention.size"))
	if err != nil {
		log.Fatal(err)
	}
	files := rotator.ListFilesWithSize()
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		diskUsedKB, err := rotator.GetDiskUsedKB(file.FullPath)
		if err != nil {
			log.Fatalf("cannot find out disk used. error=%s", err)
			return
		}
		if diskUsedKB < retentionSizeKB {
			//		fmt.Fprintf(w, "DeleteBySize(%.1fm < %.1fm): done\n", float64(diskUsedKB)/1024, float64(retentionSizeKB)/1024)
			return
		}
		if dryRun {
			//			fmt.Fprintf(w, "DeleteBySize(%.1fm > %.1fm): %s (dry run)\n", float64(diskUsedKB)/1024, float64(retentionSizeKB)/1024, file.FullPath)
		} else {
			//	fmt.Fprintf(w, "DeleteBySize(%.1fm > %.1fm): %s\n", float64(diskUsedKB)/1024, float64(retentionSizeKB)/1024, file.FullPath)
			rotator.driver.Delete(file.FullPath)
			time.Sleep(1000 * time.Millisecond) // sleep 1 second
		}
	}
}

func (rotator *Rotator) GetDiskUsedKB(path string) (int, error) {

	// todo
	// get disk used ? or under path size total?
	fileInfos, err := rotator.driver.Walk(path)
	if err != nil {
		return 0, err
	}
	var size int64
	for _, fileInfo := range fileInfos {
		size += fileInfo.Size()
	}

	return int(size), err
}
