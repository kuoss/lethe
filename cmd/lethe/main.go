package main

import (
	"log"

	"github.com/kuoss/lethe/app"
)

var Version = "development"

func main() {
	if err := app.Run(Version); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
