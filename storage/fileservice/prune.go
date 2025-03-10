package fileservice

import (
	"fmt"
	"os"
	"strings"

	"github.com/kuoss/lethe/storage/fileservice/fileutil"
)

func (s *FileService) Prune() error {
	if err := s.removeBuggyFiles(); err != nil {
		return fmt.Errorf("removeBuggyFiles err: %w", err)
	}
	if err := s.removeEmptyDirs(); err != nil {
		return fmt.Errorf("removeEmptyDirs err: %w", err)
	}
	return nil
}

func (s *FileService) removeBuggyFiles() error {
	files, err := s.driver.List("")
	if err != nil {
		return fmt.Errorf("driver.List err: %w", err)
	}

	for _, file := range files {
		if strings.HasPrefix(file, "host") || strings.HasPrefix(file, "kube") {
			if err := s.driver.Delete(file); err != nil {
				return fmt.Errorf("driver.Delete err: %w", err)
			}
		}
	}
	return nil
}

func (s *FileService) removeEmptyDirs() error {
	for _, d := range s.ListLogDirs() {
		empty, err := fileutil.IsEmpty(d.Fullpath)
		if err != nil {
			return fmt.Errorf("isEmpty err: %w", err)
		}
		if !empty {
			continue
		}
		if err := os.Remove(d.Fullpath); err != nil {
			return fmt.Errorf("os.Remove err: %w", err)
		}
	}
	return nil
}
