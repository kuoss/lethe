package main

import (
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/logger"
	"github.com/kuoss/lethe/logs"
)

var (
	Version = "unknown"
)

func main() {
	logger := logger.GetLogger()
	logger.Infof("🌊 lethe starting... version: %s", Version)
	config.LoadConfig()

	rotator := logs.NewRotator()
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	router := NewRouter()
	_ = router.Run(":3100")
}
