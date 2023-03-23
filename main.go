package main

import (
	"fmt"
	"github.com/kuoss/lethe/logs/rotator"
	"time"

	"github.com/kuoss/lethe/config"
)

func main() {
	config.LoadConfig()
	rotator := rotator.NewRotator()
	rotator.Start(time.Duration(20) * time.Minute) // 20 minutes

	fmt.Println("🌊 lethe starting...")

	r := NewRouter()
	r.Run(":3100")
}
