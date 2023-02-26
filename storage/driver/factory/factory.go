package factory

import (
	"fmt"

	"github.com/kuoss/lethe/storage/driver"
)

var driverFactories = make(map[string]StorageDriverFactory)

type StorageDriverFactory interface {
	Create(parameters map[string]interface{}) (driver.StorageDriver, error)
}

func Register(name string, factory StorageDriverFactory) {
	if factory == nil {
		panic("Must not provide nil StorageDriverFactory")
	}
	_, registered := driverFactories[name]
	if registered {
		panic(fmt.Sprintf("StorageDriverFactory named %s already registered", name))
	}
	driverFactories[name] = factory
}

func Get(name string, parameters map[string]interface{}) (driver.StorageDriver, error) {
	driverFactory, ok := driverFactories[name]
	if !ok {
		return nil, fmt.Errorf("invalid StorageDriver named %s", name)
	}
	return driverFactory.Create(parameters)
}
