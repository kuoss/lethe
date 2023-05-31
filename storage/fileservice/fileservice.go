package fileservice

import (
	"fmt"
	"os"

	"github.com/kuoss/lethe/config"
	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"
)

type FileService struct {
	config *config.Config
	driver storagedriver.Driver
}

func New(cfg *config.Config) (*FileService, error) {
	driver, err := factory.Get("filesystem", map[string]interface{}{"RootDirectory": cfg.LogDataPath()})
	if err != nil {
		return nil, fmt.Errorf("factory.Get err: %w", err)
	}
	// TODO: use driver
	err = os.MkdirAll(cfg.LogDataPath(), 0755)
	if err != nil {
		return nil, fmt.Errorf("mkdirAll err: %w", err)
	}
	return &FileService{cfg, driver}, nil
}

func (s *FileService) Config() *config.Config {
	return s.config
}
