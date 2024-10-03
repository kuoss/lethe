package main

import (
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/rotator"
	"github.com/kuoss/lethe/router"
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
		logger.Fatalf("new config err: %s", err.Error())
	}
	logger.Infof("%+v", cfg)

	// services
	fileService, err := fileservice.New(cfg)
	if err != nil {
		logger.Fatalf("new fileservice err: %s", err.Error())
	}
	logService := logservice.New(fileService)
	queryService := queryservice.New(logService)

	// start rotator
	rotator := rotator.New(cfg, fileService)
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	// run router
	h := router.New(cfg, fileService, queryService)
	err = h.Run()
	if err != nil {
		logger.Fatalf("router run err: %s", err.Error())
	}
}
