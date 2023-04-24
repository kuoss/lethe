package main

import (
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/logs/rotator"
)

var (
	Version = "unknown"
)

func main() {
	logger.Infof("🌊 lethe starting... version: %s", Version)
	err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("error on LoadConfig: %s", err)
	}

	rotator := rotator.NewRotator()
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	router := NewRouter()
	_ = router.Run(":3100")
}
