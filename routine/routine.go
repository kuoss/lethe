package routine

import (
	"log"
	"time"

	"github.com/kuoss/lethe/file"
)

func Start(interval time.Duration) {
	go routineLoop(interval)
}

func routineLoop(interval time.Duration) {
	for {
		file.DeleteByAge(false)
		file.DeleteBySize(false)
		file.Cleansing()

		log.Printf("routineLoop... sleep %s\n", interval)
		time.Sleep(interval)
	}
}
