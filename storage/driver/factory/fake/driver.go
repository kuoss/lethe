package fake

import (
	"io"

	storagedriver "github.com/kuoss/lethe/storage/driver"
)

const (
	driverName = "fake"
)

type driver struct{}

func New() storagedriver.Driver {
	return &driver{}
}
func (*driver) Name() string {
	return driverName
}
func (*driver) GetContent(string) ([]byte, error) {
	return []byte{}, nil
}
func (*driver) PutContent(string, []byte) error {
	return nil
}
func (*driver) Reader(string) (io.ReadCloser, error) {
	return &io.PipeReader{}, nil
}
func (*driver) Stat(string) (storagedriver.FileInfo, error) {
	return nil, nil
}
func (*driver) List(string) ([]string, error) {
	return []string{}, nil
}
func (*driver) Move(string, string) error {
	return nil
}
func (*driver) Delete(s string) error {
	return nil
}
func (*driver) Walk(string) ([]storagedriver.FileInfo, error) {
	return nil, nil
}
func (*driver) WalkDir(string) ([]string, error) {
	return []string{}, nil
}
func (*driver) Writer(string) (storagedriver.FileWriter, error) {
	return nil, nil
}
func (*driver) RootDirectory() string {
	return ""
}
