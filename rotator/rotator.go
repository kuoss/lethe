package rotator

import (
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/fileservice"
)

type Rotator struct {
	interval    time.Duration
	fileService *fileservice.FileService
}

func New(cfg *config.Config, fileService *fileservice.FileService) *Rotator {
	return &Rotator{cfg.Rotator.Interval, fileService}
}

func (r *Rotator) Start() {
	go r.routineLoop()
}

func (r *Rotator) routineLoop() {
	for {
		r.RunOnce()
		logger.Infof("routineLoop... sleep %s", r.interval)
		time.Sleep(r.interval)
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
