package app

import (
	"fmt"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/rotator"
	"github.com/kuoss/lethe/router"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice"
	"github.com/kuoss/lethe/storage/queryservice"
)

type IApp interface {
	Run(version string) error
}

type App struct{}

func (a App) Run(version string) error {
	logger.Infof("Starting Lethe 💧 version=%s", version)

	// Load configuration
	cfg, err := config.New(version)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	logger.Infof("Loaded configuration: %v", cfg)

	// Initialize services
	fileService, err := fileservice.New(cfg)
	if err != nil {
		return fmt.Errorf("new fileservice err: %w", err)
	}
	logService := logservice.New(cfg, fileService)
	queryService := queryservice.New(cfg, logService)

	// Start rotater
	myRotator := rotator.New(cfg, fileService)
	myRotator.Start()

	// Run router
	myRouter := router.New(cfg, fileService, queryService)
	return myRouter.Run()
}
