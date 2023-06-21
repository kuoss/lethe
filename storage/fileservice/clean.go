package fileservice

import (
	"fmt"

	"github.com/kuoss/common/logger"
)

func (s *FileService) Clean() {
	s.removeFilesWithPrefix("host")
	s.removeFilesWithPrefix("kube")
}

func (s *FileService) removeFilesWithPrefix(prefix string) {
	files, err := s.driver.Walk(fmt.Sprintf("%s/%s.*", s.config.LogDataPath(), prefix))
	if err != nil {
		logger.Warnf("glob err: %s, prefix: %s", err.Error(), prefix)
		return
	}
	if len(files) < 1 {
		return
	}
	logger.Warnf("cleansing files prefix: %s", prefix)
	for _, file := range files {
		logger.Infof("remove file: %s", file)
		err := s.driver.Delete(file.Fullpath())
		if err != nil {
			logger.Warnf("remove err: %s, file: %s", err.Error(), file)
		}
	}
}
