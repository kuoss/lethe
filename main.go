package main

import (
	"fmt"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/logs"
)

func main() {
	config.LoadConfig()
	rotator := logs.NewRotator()
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	fmt.Println("🌊 lethe starting...")

	r := NewRouter()
	_ = r.Run()
}
