package logs

import (
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"
	"log"
	"time"
)

type Rotator struct {
	driver driver.StorageDriver
}

func NewRotator() *Rotator {
	d, _ := factory.Get(config.GetConfig().GetString("storage.driver"), map[string]interface{}{"RootDirectory": config.GetLogRoot()})
	return &Rotator{driver: d}
}

func (rotator *Rotator) Start(interval time.Duration) {
	go rotator.routineLoop(interval)
}

func (rotator *Rotator) routineLoop(interval time.Duration) {

	for {
		rotator.DeleteByAge(true)
		rotator.DeleteBySize(true)
		rotator.Cleansing()

		log.Printf("routineLoop... sleep %s\n", interval)
		time.Sleep(interval)
	}
}
