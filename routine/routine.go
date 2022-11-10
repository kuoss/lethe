package routine

import (
	"log"
	"time"

	"github.com/kuoss/lethe/file"
)

func Start(interval time.Duration) {
	log.Printf("Routine started... interval=%s", interval)
	go routineLoop(interval)
}

func routineLoop(interval time.Duration) {
	for {
		TaskDelete()
		time.Sleep(interval)
	}
}

func TaskDelete() {
	log.Println("TaskDelete started...")
	file.DeleteByAge(false)
	file.DeleteBySize(false)
}
