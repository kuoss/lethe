package fileservice

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kuoss/common/logger"
)

func (s *FileService) Prune() {
	warns, err := s.removeUnneccesaryFiles()
	if err != nil {
		logger.Errorf("removeUnneccesaryFiles err: %v", err)
	}
	if len(warns) > 0 {
		logger.Warnf("removeUnneccesaryFiles warns: %v", warns)
	}

	warns, err = s.removeOldEmptyDirs()
	if err != nil {
		logger.Errorf("removeOldEmptyDirs err: %v", err)
	}
	if len(warns) > 0 {
		logger.Warnf("removeOldEmptyDirs warns: %v", warns)
	}
}

func (s *FileService) removeUnneccesaryFiles() ([]string, error) {
	files, err := s.driver.List("")
	if err != nil {
		return nil, fmt.Errorf("driver.List err: %w", err)
	}

	warns := []string{}
	for _, file := range files {
		if strings.HasPrefix(file, "host") || strings.HasPrefix(file, "kube") {
			if err := s.driver.Delete(file); err != nil {
				warns = append(warns, fmt.Sprintf("failed to delete '%s': %v", file, err))
			}
		}
	}
	return warns, nil
}

func (s *FileService) removeOldEmptyDirs() ([]string, error) {
	warns := []string{}
	for _, d := range s.ListLogDirs() {
		oldEmpty, err := isOldEmptyDir(d.Fullpath)
		if err != nil {
			warns = append(warns, fmt.Sprintf("isOldEmptyDir check failed for '%s': %v", d.Fullpath, err))
			continue
		}
		if !oldEmpty {
			continue
		}
		if err := os.Remove(d.Fullpath); err != nil {
			warns = append(warns, fmt.Sprintf("failed to remove old empty dir '%s': %v", d.Fullpath, err))
		}
	}
	return warns, nil
}

func isOldEmptyDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("path is not a directory")
	}

	dayOld := time.Now().Add(-24 * time.Hour)
	if info.ModTime().After(dayOld) {
		return false, nil
	}

	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	entries, err := dir.Readdirnames(1)
	if err != nil {
		if err == io.EOF {
			return true, nil
		}
		return false, err
	}

	return len(entries) == 0, nil
}
