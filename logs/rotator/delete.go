package rotator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
)

// DELETE
func (rotator *Rotator) DeleteByAge() error {
	retentionTime := config.Viper().GetString("retention.time")
	duration, err := util.GetDurationFromAge(retentionTime)
	if err != nil {
		return fmt.Errorf("error on GetDurationFromAge: %w", err)
	}
	point := strings.Replace(clock.GetNow().Add(-duration).UTC().String()[0:13], " ", "_", 1)
	files, err := rotator.ListFiles()
	if err != nil {
		return fmt.Errorf("error on ListFiles: %w", err)
	}
	if len(files) < 1 {
		logger.Infof("DeleteByAge( < %s): no files. done.", point)
		return nil
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		if file.Name < point {
			logger.Infof("DeleteByAge(%s < %s): %s", file.Name, point, file.FullPath)
			err := rotator.driver.Delete(file.FullPath)
			if err != nil {
				logger.Errorf("error on Delete: %s", err)
				continue
			}
		}
	}
	logger.Infof("DeleteByAge(%s): Done\n", point)
	return nil
}

func (rotator *Rotator) DeleteBySize() error {
	retentionSizeBytes, err := util.StringToBytes(config.Viper().GetString("retention.size"))
	if err != nil {
		return fmt.Errorf("error on StringToBytes: %w", err)
	}
	files, err := rotator.ListFiles()
	if err != nil {
		return fmt.Errorf("error on ListFiles: %w", err)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		usedBytes, err := rotator.GetUsedBytes(config.GetLogDataPath())
		if err != nil {
			return fmt.Errorf("error on GetUsedBytes: %w", err)
		}
		if usedBytes < retentionSizeBytes {
			logger.Infof("DeleteBySize(%d < %d): DONE", usedBytes, retentionSizeBytes)
			return nil
		}
		logger.Infof("DeleteBySize(%d > %d): %s", usedBytes, retentionSizeBytes, file.FullPath)
		err = rotator.driver.Delete(file.FullPath)
		if err != nil {
			logger.Errorf("error on Delte: %s", err)
		}
	}
	return nil
}
