package factory

import (
	"fmt"

	storagedriver "github.com/kuoss/lethe/storage/driver"
)

var driverFactories = make(map[string]StorageDriverFactory)

type StorageDriverFactory interface {
	Create(parameters map[string]interface{}) (storagedriver.Driver, error)
}

func Register(name string, factory StorageDriverFactory) error {
	if factory == nil {
		return fmt.Errorf("factory is nil")
	}
	if _, exist := driverFactories[name]; exist {
		return fmt.Errorf("factory name duplicated: %s", name)
	}
	driverFactories[name] = factory
	return nil
}

func Get(name string, parameters map[string]interface{}) (storagedriver.Driver, error) {
	driverFactory, ok := driverFactories[name]
	if !ok {
		return nil, fmt.Errorf("invalid driver name: %s", name)
	}
	return driverFactory.Create(parameters)
}
