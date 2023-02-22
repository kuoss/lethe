package filesystem

import storagedriver "github.com/kuoss/lethe/storage/driver"

type filesystemDriverFactory struct{}

func (factory *filesystemDriverFactory) Create(parameters map[string]interface{}) (storagedriver.StorageDriver, error) {
	return New(DriverParameters{RootDirectory: parameters["RootDirectory"].(string)}), nil
}
