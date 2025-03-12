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
	return &Rotator{
		interval:    cfg.Retention.RotationInterval,
		fileService: fileService,
	}
}

func (r *Rotator) Start() {
	go r.rotatePeriodically()
}

func (r *Rotator) rotatePeriodically() {
	for {
		r.Rotate()
		logger.Infof("sleep %s", r.interval)
		time.Sleep(r.interval)
	}
}

func (r *Rotator) Rotate() {
	if err := r.fileService.DeleteByAge(); err != nil {
		logger.Errorf("DeleteByAge err: %v", err)
	}
	if err := r.fileService.DeleteBySize(); err != nil {
		logger.Errorf("DeleteBySize err: %v", err)
	}
	r.fileService.Prune()
}
