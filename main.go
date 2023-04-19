package main

import (
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/logger"
	"github.com/kuoss/lethe/logs/rotator"
)

var (
	Version = "unknown"
)

func main() {
	logger.Infof("🌊 lethe starting... version: %s", Version)
	config.LoadConfig()

	rotator := rotator.NewRotator()
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	router := NewRouter()
	_ = router.Run(":3100")
}
