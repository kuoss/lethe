package rotator

import (
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
)

type Rotator struct {
	config      *config.Config
	fileService *fileservice.FileService
}

func New(cfg *config.Config, fileService *fileservice.FileService) *Rotator {
	return &Rotator{cfg, fileService}
}

func (r *Rotator) Start(interval time.Duration) {
	go r.routineLoop(interval)
}

func (r *Rotator) routineLoop(interval time.Duration) {
	for {
		r.RunOnce()
		logger.Infof("routineLoop... sleep %s", interval)
		time.Sleep(interval)
	}
}

func (r *Rotator) RunOnce() {
	err := r.fileService.DeleteByAge()
	if err != nil {
		logger.Errorf("deleteByAge err: %s", err.Error())
	}
	err = r.fileService.DeleteBySize()
	if err != nil {
		logger.Errorf("deleteBySize err: %s", err.Error())
	}
	r.fileService.Clean()
}
