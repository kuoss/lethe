package main

import (
	"fmt"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/logs"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"time"
)

func main() {
	config.LoadConfig()

	rotator := logs.NewRotator()
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	fmt.Println("🌊 lethe starting...")

	// start server
	r := NewRouter()
	r.Run()
}
