package filesystem

import (
	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"
)

type filesystemDriverFactory struct{}

func (factory *filesystemDriverFactory) Create(params map[string]interface{}) (storagedriver.Driver, error) {
	return New(Params{RootDirectory: params["RootDirectory"].(string)}), nil
}

func init() {
	err := factory.Register(driverName, &filesystemDriverFactory{})
	if err != nil {
		panic(err)
	}
}
