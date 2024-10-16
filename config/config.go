package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/units"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/spf13/viper"

	"github.com/kuoss/common/logger"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
)

type Config struct {
	Version string

	Query     Query
	Retention Retention
	Storage   Storage
	Web       Web
}

type Query struct {
	Limit   int
	Timeout time.Duration
}

type Retention struct {
	Size             int64
	Time             time.Duration
	SizeStrategy     string
	RotationInterval time.Duration
}

type Storage struct {
	LogDataPath string
}

type Web struct {
	ListenAddress string
	GinMode       string
}

func New(version string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("lethe")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(".", "etc"))
	v.AddConfigPath(filepath.Join("..", "etc"))

	// Set default values
	v.SetDefault("query.limit", 1000)
	v.SetDefault("query.timeout", "20s")
	v.SetDefault("retention.size", 0)
	v.SetDefault("retention.time", "15d")
	v.SetDefault("retention.size_strategy", "file")
	v.SetDefault("retention.rotation_interval", "20s")
	v.SetDefault("storage.log_data_path", "data/log")
	v.SetDefault("web.listen_address", ":6060")
	v.SetDefault("web.gin_mode", gin.ReleaseMode) // "release"

	// Read the configuration file, if it exists
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			_ = showFileContent(v.ConfigFileUsed())
			// Config file was found but another error occurred
			return nil, fmt.Errorf("Config file has error: %w", err)
		}
		logger.Warnf("Config file not found; continuing with defaults")
	} else {
		_ = showFileContent(v.ConfigFileUsed())
	}

	// units
	retentionSizeString := v.GetString("retention.size")
	retentionSize, err := parseSize(retentionSizeString)
	if err != nil {
		return nil, fmt.Errorf("parse retention size err: %w, size: %s", err, retentionSizeString)
	}

	retentionTimeString := v.GetString("retention.time")
	retentionTime, err := model.ParseDuration(retentionTimeString)
	if err != nil {
		return nil, fmt.Errorf("parse retention time err: %w, time: %s", err, retentionTimeString)
	}

	return &Config{
		Version: version,
		Query: Query{
			Limit:   v.GetInt("query.limit"),
			Timeout: v.GetDuration("query.timeout"),
		},
		Retention: Retention{
			Size:             int64(retentionSize),
			Time:             time.Duration(retentionTime),
			SizeStrategy:     v.GetString("retention.size_strategy"),
			RotationInterval: v.GetDuration("retention.rotation_interval"),
		},
		Storage: Storage{
			LogDataPath: v.GetString("storage.log_data_path"),
		},
		Web: Web{
			ListenAddress: v.GetString("web.listen_address"),
			GinMode:       v.GetString("web.gin_mode"),
		},
	}, nil
}

func showFileContent(filePath string) error {
	logger.Infof("Read config file: %s", filePath)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file err: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

func parseSize(sizeString string) (int64, error) {
	// Handle legacy units and ensure the size string is uppercase for parsing
	unitReplacements := map[string]string{
		"k": "KB",
		"m": "MB",
		"g": "GB",
	}

	// Replace legacy units
	for oldUnit, newUnit := range unitReplacements {
		if strings.HasSuffix(sizeString, oldUnit) {
			sizeString = strings.ReplaceAll(sizeString, oldUnit, newUnit)
		}
	}

	// Parse the size string using units package
	size, err := units.ParseBase2Bytes(sizeString)
	if err != nil {
		return 0, fmt.Errorf("failed to parse size '%s': %w", sizeString, err)
	}
	return int64(size), nil
}
