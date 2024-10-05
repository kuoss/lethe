package main

import (
	"os"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/app"
)

var (
	Version = "development"

	myApp app.IApp
	exit  = os.Exit
)

func main() {
	logger.Infof("Starting Lethe, version: %s", Version)
	if err := myApp.New(Version); err != nil {
		logger.Errorf("Failed to create app: %s", err.Error())
		exit(1)
	}
	if err := myApp.Run(); err != nil {
		logger.Errorf("Failed to run app: %s", err.Error())
		exit(1)
	}
}
