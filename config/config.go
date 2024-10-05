package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/common/model"
	"github.com/spf13/viper"

	"github.com/kuoss/common/logger"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Config struct {
	Version string

	Query     Query
	Retention Retention
	Rotator   Rotator
	Storage   Storage
	Web       Web
}

type Query struct {
	Limit   int
	Timeout time.Duration
}

type Retention struct {
	SizingStrategy string
	Size           int64
	Time           time.Duration
}

type Rotator struct {
	Interval time.Duration
}

type Storage struct {
	LogDataPath string
}

type Web struct {
	ListenAddress string
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
	v.SetDefault("retention.size", "1000Gi")
	v.SetDefault("retention.time", "15d")
	v.SetDefault("retention.sizingStrategy", "file")
	v.SetDefault("rotator.interval", "20s")
	v.SetDefault("storage.logDataPath", "/data/log")
	v.SetDefault("web.listenAddress", ":6060")

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
	retentionSize, err := resource.ParseQuantity(retentionSizeString)
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
			Size:           retentionSize.Value(),
			Time:           time.Duration(retentionTime),
			SizingStrategy: v.GetString("retention.sizingStrategy"),
		},
		Rotator: Rotator{
			Interval: v.GetDuration("rotator.interval"),
		},
		Storage: Storage{
			LogDataPath: v.GetString("storage.logDataPath"),
		},
		Web: Web{
			ListenAddress: v.GetString("web.listenAddress"),
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
