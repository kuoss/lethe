package rotator

import (
	"log"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"
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
		rotator.RunOnce()
		log.Printf("routineLoop... sleep %s\n", interval)
		time.Sleep(interval)
	}
}

func (rotator *Rotator) RunOnce() {
	var err error

	err = rotator.DeleteByAge()
	if err != nil {
		logger.Errorf("error on DeleteByAge: %s", err)
	}

	err = rotator.DeleteBySize()
	if err != nil {
		logger.Errorf("error on DeleteBySize: %s", err)
	}

	rotator.Cleansing()
}
