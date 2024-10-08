package fileservice

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/clock"
)

func (s *FileService) DeleteByAge() error {
	retentionTime := s.Config.Retention.Time
	if retentionTime == 0 {
		logger.Infof("retentionTime is 0 (DeleteByAge skipped)")
		return nil
	}
	point := strings.Replace(clock.Now().Add(-retentionTime).UTC().String()[0:13], " ", "_", 1)
	files, err := s.ListFiles()
	if err != nil {
		return fmt.Errorf("listFiles err: %w", err)
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
			logger.Infof("DeleteByAge(%s < %s): %s", file.Name, point, file.Fullpath)
			err := s.driver.Delete(file.Subpath)
			if err != nil {
				logger.Errorf("delete err: %s", err.Error())
				continue
			}
		}
	}
	logger.Infof("DeleteByAge(%s): DONE", point)
	return nil
}

func (s *FileService) DeleteBySize() error {
	retentionSize := s.Config.Retention.Size
	if retentionSize == 0 {
		logger.Infof("retentionSize is 0 (DeleteBySize skipped)")
		return nil
	}

	// ListFiles use driver.Walk
	files, err := s.ListFiles()
	if err != nil {
		return fmt.Errorf("listFiles err: %w", err)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	// calculate sum of all files size
	var used int64 = 0
	for _, file := range files {
		used += file.Size
	}

	var deleteSize int64 = 0
	var deleteFiles []LogFile
	for _, file := range files {
		if used-deleteSize < retentionSize {
			break
		}
		deleteFiles = append(deleteFiles, file)
		deleteSize += file.Size
	}
	logger.Infof("DeleteBySize: try to flush %d files, %d bytes", len(deleteFiles), deleteSize)
	for _, file := range deleteFiles {
		err := s.driver.Delete(file.Subpath)
		if err != nil {
			logger.Errorf("delete err: %s", err.Error())
		}
		logger.Infof("DeleteBySize(%d > %d): %s(%d)", used, retentionSize, file.Fullpath, file.Size)
	}
	logger.Infof("DeleteBySize(%d < %d): DONE", used, retentionSize)
	return nil
}
