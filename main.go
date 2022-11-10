package main

import (
	"log"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/routine"
)

func main() {
	config.LoadConfig()
	routine.Start(time.Duration(20) * time.Minute) // 20 minutes
	log.Printf("🌊 lethe starting...")

	// start server
	r := NewRouter()
	r.Run()
}

type Query struct {
	Expr string `form:"expr"`
}
