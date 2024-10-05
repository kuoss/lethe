package fileservice

import (
	"fmt"
	"os"

	"github.com/kuoss/lethe/config"
	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"
)

type FileService struct {
	Config *config.Config
	driver storagedriver.Driver
}

func New(cfg *config.Config) (*FileService, error) {
	driver, err := factory.Get("filesystem", map[string]any{"RootDirectory": cfg.Storage.LogDataPath})
	if err != nil {
		return nil, fmt.Errorf("factory.Get err: %w", err)
	}

	if err := os.MkdirAll(cfg.Storage.LogDataPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("os.MkdirAll err: %w", err)
	}

	return &FileService{cfg, driver}, nil
}
