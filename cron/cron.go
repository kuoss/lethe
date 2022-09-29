package cron

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kuoss/lethe/file"
)

const (
	TaskDelete = 1
)

type Job struct {
	Task     int
	Interval int
}

func Start(jobs []Job) {
	s := gocron.NewScheduler(time.UTC)
	for _, job := range jobs {
		switch job.Task {
		case TaskDelete:
			s.Every(job.Interval).Seconds().Do(DoTaskDelete)
		}
	}
	s.StartAsync()
}

func DoTaskDelete() {
	log.Println("TaskDelete started...")
	file.DeleteByAge(false)
	file.DeleteBySize(false)
}
