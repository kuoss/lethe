package main

import (
	"os"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/app"
)

var (
	Version = "development" // Version will be overwritten by ldflags

	ap   app.IApp = app.App{}
	exit          = os.Exit
)

func main() {
	if err := ap.Run(Version); err != nil {
		logger.Errorf("Failed to run app: %v", err)
		exit(1)
	} else {
		exit(0)
	}
}
