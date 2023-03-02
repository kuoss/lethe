package logs

import (
	"github.com/kuoss/lethe/config"
)

func (rotator *Rotator) GetUsedBytes(path string) (int, error) {
	if config.GetConfig().GetString("retention.sizingStrategy") == "disk" {
		return rotator.GetDiskUsedBytes(path)
	}
	return rotator.GetFilesUsedBytes(path)
}

func (rotator *Rotator) GetFilesUsedBytes(path string) (int, error) {
	fileInfos, err := rotator.driver.Walk(path)
	if err != nil {
		return 0, err
	}
	var size int
	for _, fileInfo := range fileInfos {
		size += int(fileInfo.Size())
	}
	return size, err
}
