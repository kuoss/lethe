package fileservice

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/clock"
	"golang.org/x/sys/unix"
)

func (s *FileService) DeleteByAge() error {
	retentionTime := s.config.RetentionTime()
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
	retentionSize := s.config.RetentionSize()
	if retentionSize == 0 {
		logger.Infof("retentionSize is 0 (DeleteBySize skipped)")
		return nil
	}
	files, err := s.ListFiles()
	if err != nil {
		return fmt.Errorf("listFiles err: %w", err)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	for _, file := range files {
		usedBytes, err := s.GetUsedBytes(".")
		if err != nil {
			return fmt.Errorf("getUsedBytes err: %w", err)
		}
		if usedBytes < retentionSize {
			logger.Infof("DeleteBySize(%d < %d): DONE", usedBytes, retentionSize)
			return nil
		}
		logger.Infof("DeleteBySize(%d > %d): %s", usedBytes, retentionSize, file.Fullpath)
		err = s.driver.Delete(file.Subpath)
		if err != nil {
			logger.Errorf("delete err: %s", err.Error())
		}
	}
	return nil
}

func (s *FileService) GetUsedBytes(subpath string) (int, error) {
	if s.config.RetentionSizingStrategy() == "disk" {
		return s.GetUsedBytesFromDisk(subpath)
	}
	return s.GetUsedBytesFromFiles(subpath)
}

func (s *FileService) GetUsedBytesFromFiles(subpath string) (int, error) {
	fileInfos, err := s.driver.Walk(subpath)
	if err != nil {
		return 0, err
	}
	var size int
	for _, fileInfo := range fileInfos {
		size += int(fileInfo.Size())
	}
	return size, err
}

func (s *FileService) GetUsedBytesFromDisk(path string) (int, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}
	return int(stat.Blocks - stat.Bavail*uint64(stat.Bsize)), nil
}
