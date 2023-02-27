package logs

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

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
		diskUsedBytes, err := rotator.GetDiskUsedBytes(config.GetLogRoot())
		if err != nil {
			log.Fatalf("Cannot get disk used bytes: %s", err)
			return
		}
		if diskUsedBytes < retentionSizeBytes {
			fmt.Printf("DeleteBySize(%d < %d): Done\n", diskUsedBytes, retentionSizeBytes)
			return
		}
		fmt.Printf("DeleteBySize(%d > %d): %s\n", diskUsedBytes, retentionSizeBytes, file.FullPath)
		rotator.driver.Delete(file.FullPath)
		// time.Sleep(500 * time.Millisecond)
	}
}

func (rotator *Rotator) GetDiskUsedBytes(path string) (int, error) {
	if os.Getenv("TEST_MODE") == "1" {
		return rotator.GetDiskUsedBytesInTest(path)
	}
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return -1, err
	}
	return int(int64(stat.Blocks-stat.Bavail) * stat.Bsize), nil
}

func (rotator *Rotator) GetDiskUsedBytesInTest(path string) (int, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return -1, err
	}
	avail := int(int64(stat.Bavail) * stat.Bsize)
	initialAvail, err := strconv.Atoi(os.Getenv("TEST_INITIAL_DISK_AVAILABLE_BYTES"))
	if err != nil {
		return -1, err
	}
	return initialAvail - avail, nil
}

func (rotator *Rotator) GetFilesUsedBytes(path string) (int, error) {
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
