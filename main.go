package main

import (
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/handler"
	"github.com/kuoss/lethe/rotator"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/queryservice"
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

	// services
	fileService, err := fileservice.New(cfg)
	if err != nil {
		logger.Fatalf("fileSvc.New err: %s", err.Error())
	}
	logService := logservice.New(fileService)
	queryService := queryservice.New(logService)

	// start rotator
	rotator := rotator.New(cfg, fileService)
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	// run router
	h := handler.New(cfg, fileService, queryService)
	err = h.Run()
	if err != nil {
		logger.Fatalf("run err: %s", err.Error())
	}
}
