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
	logger.Infof("💧 Lethe starting... version: %s", Version)

	// Load configuration
	cfg := mustConfig(Version)

	// Initialize services
	fileService := mustFileService(cfg)
	logService := logservice.New(fileService)
	queryService := queryservice.New(logService)

	// Start the rotator
	startRotator(cfg, fileService)

	// Run the router
	runRouter(cfg, fileService, queryService)
}

func mustConfig(version string) *config.Config {
	cfg, err := config.New(version)
	if err != nil {
		logger.Fatalf("Failed to create new config: %s", err.Error())
	}
	logger.Infof("Loaded configuration: %+v", cfg)
	return cfg
}

func mustFileService(cfg *config.Config) *fileservice.FileService {
	fileService, err := fileservice.New(cfg)
	if err != nil {
		logger.Fatalf("Failed to create new file service: %s", err.Error())
	}
	return fileService
}

func startRotator(cfg *config.Config, fileService *fileservice.FileService) {
	rotator := rotator.New(cfg, fileService)
	rotator.Start(20 * time.Minute)
	logger.Infof("Log rotator started with interval: %v", 20*time.Minute)
}

func runRouter(cfg *config.Config, fileService *fileservice.FileService, queryService *queryservice.QueryService) {
	r := router.New(cfg, fileService, queryService)
	if err := r.Run(); err != nil {
		logger.Fatalf("Failed to run router: %s", err.Error())
	}
}
