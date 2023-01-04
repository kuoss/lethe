package main

import (
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/logs"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"log"
	"time"
)

func main() {
	config.LoadConfig()

	rotator := logs.NewRotator()
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	log.Println("🌊 lethe starting...")

	// start server
	r := NewRouter()
	r.Run()
}
