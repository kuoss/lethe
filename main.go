package main

import (
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/rotator"
	"github.com/kuoss/lethe/router"
	"github.com/kuoss/lethe/storage/fileservice"
)

var (
	Version = "development"
)

func main() {

	logger.Infof("🌊 lethe starting... version: %s", Version)

	// config
	cfg, err := config.New(Version)
	if err != nil {
		logger.Fatalf("config.New err: %s", err.Error())
	}

	// fileService
	fileService, err := fileservice.New(cfg)
	if err != nil {
		logger.Fatalf("fileSvc.New err: %s", err.Error())
	}

	// start rotator
	rotator := rotator.New(cfg, fileService)
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	// run router
	router := router.New(fileService)
	err = router.Run(cfg.WebListenAddress())
	if err != nil {
		logger.Fatalf("run err: %s", err.Error())
	}
}
