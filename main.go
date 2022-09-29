package main

import (
	"log"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/cron"
)

func main() {
	config.LoadConfig()
	cron.Start([]cron.Job{
		{Task: cron.TaskDelete, Interval: 20 * 60}, // 20 minutes
	})
	log.Printf("🌊 lethe starting...")
	startServer()
}

type Query struct {
	Expr string `form:"expr"`
}

func startServer() {
	r := NewRouter()
	r.Run()
}
