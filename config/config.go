package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/kuoss/common/logger"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/kuoss/lethe/util"
	"github.com/spf13/viper"
)

type Config struct {
	Version string

	limit                   int
	logDataPath             string
	retentionSize           int
	retentionTime           time.Duration
	retentionSizingStrategy string
	timeout                 time.Duration
	webListenAddress        string
}

func New(version string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("lethe")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(".", "etc"))
	v.AddConfigPath(filepath.Join("..", "etc"))

	// Set default values
	v.SetDefault("version", "test")
	v.SetDefault("limit", 1000)
	v.SetDefault("storage.log_data_path", "/var/data/log")
	v.SetDefault("retention.sizingStrategy", "file")
	v.SetDefault("web.listen_address", ":6060")
	v.SetDefault("timeout", 20*time.Second)
	v.SetDefault("retention.time", "15d")
	v.SetDefault("retention.size", "1000g")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warnf("Configuration lethe.yaml is not provided\n")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("readInConfig err: %w", err)
		}
	}

	retentionSize, err := util.StringToBytes(v.GetString("retention.size"))
	if err != nil {
		return nil, fmt.Errorf("stringToBytes err: %w", err)
	}

	retentionTimeString := v.GetString("retention.time")
	retentionTime, err := util.GetDurationFromAge(retentionTimeString)
	if err != nil {
		return nil, fmt.Errorf("getDurationFromAge err: %w", err)
	}

	return &Config{
		Version:                 version,
		limit:                   1000,
		logDataPath:             v.GetString("storage.log_data_path"),
		retentionSize:           retentionSize, // 1000g
		retentionTime:           retentionTime, // 15d
		retentionSizingStrategy: v.GetString("retention.sizingStrategy"),
		timeout:                 v.GetDuration("timeout"),
		webListenAddress:        v.GetString("web.listen_address"),
	}, nil
}

func (c *Config) Limit() int {
	return c.limit
}

func (c *Config) LogDataPath() string {
	return c.logDataPath
}
func (c *Config) SetLogDataPath(logDataPath string) {
	c.logDataPath = logDataPath
}

func (c *Config) RetentionSize() int {
	return c.retentionSize
}
func (c *Config) SetRetentionSize(retentionSize int) {
	c.retentionSize = retentionSize
}

func (c *Config) RetentionTime() time.Duration {
	return c.retentionTime
}
func (c *Config) SetRetentionTime(retentionTime time.Duration) {
	c.retentionTime = retentionTime
}

func (c *Config) RetentionSizingStrategy() string {
	return c.retentionSizingStrategy
}

func (c *Config) Timeout() time.Duration {
	return c.timeout
}

func (c *Config) WebListenAddress() string {
	return c.webListenAddress
}
