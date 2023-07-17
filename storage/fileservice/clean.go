package fileservice

import (
	"strings"

	"github.com/kuoss/common/logger"
)

func (s *FileService) Clean() {
	s.removeFilesWithPrefix("host")
	s.removeFilesWithPrefix("kube")
}

func (s *FileService) removeFilesWithPrefix(prefix string) {
	files, err := s.driver.List("")
	if err != nil {
		logger.Warnf("list err: %s, prefix: %s", err.Error(), prefix)
		return
	}
	if len(files) < 1 {
		return
	}
	logger.Warnf("cleanning files prefix: %s", prefix)
	for _, file := range files {
		if strings.HasPrefix(file, prefix) {
			logger.Infof("remove file: %s", file)
			err := s.driver.Delete(file)
			if err != nil {
				logger.Warnf("remove err: %s, file: %s", err.Error(), file)
			}
		}
	}
}
