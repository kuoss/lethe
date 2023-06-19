package config

import (
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/kuoss/lethe/util"
	"github.com/spf13/viper"
)

type Config struct {
	limit                   int
	logDataPath             string
	retentionSize           int
	retentionTime           time.Duration
	retentionSizingStrategy string
	timeout                 time.Duration
	version                 string
	webListenAddress        string
}

func New(version string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("lethe")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(".", "etc"))
	v.AddConfigPath(filepath.Join("..", "etc"))

	// Set default values
	v.SetDefault("storage.log_data_path", "./tmp/log")
	v.SetDefault("retention.sizingStrategy", "file")
	v.SetDefault("web.listen_address", ":6060")
	v.SetDefault("timeout", 20*time.Second)
	// fixme: decide default value
	v.SetDefault("retention.time", "15d")
	v.SetDefault("retention.size", "1000g")

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found;
			return &Config{
				limit:                   1000,
				logDataPath:             "./tmp/log",
				retentionSize:           100 * 1024 * 1024,   // fixme: we should decide default value
				retentionTime:           15 * 24 * time.Hour, // fixme: we should decide default value
				retentionSizingStrategy: "file",
				timeout:                 20 * time.Second,
				version:                 version,
				webListenAddress:        ":6060",
			}, nil
		} else {
			// Config file was found but another error was produced
			return &Config{}, fmt.Errorf("readInConfig err: %w", err)
		}
	}

	err = v.Unmarshal(&v)
	if err != nil {
		return &Config{}, fmt.Errorf("unmarshal err: %w", err)
	}

	retentionSize, err := util.StringToBytes(v.GetString("retention.size"))
	if err != nil {
		return &Config{}, fmt.Errorf("stringToBytes err: %w", err)
	}

	retentionTimeString := v.GetString("retention.time")
	retentionTime, err := util.GetDurationFromAge(retentionTimeString)
	if err != nil {
		return &Config{}, fmt.Errorf("getDurationFromAge err: %w", err)
	}

	return &Config{
		limit:                   1000,
		logDataPath:             v.GetString("storage.log_data_path"),
		retentionSize:           retentionSize, // fixme: we should decide default value
		retentionTime:           retentionTime, // fixme: we should decide default value
		retentionSizingStrategy: v.GetString("retention.sizingStrategy"),
		timeout:                 v.GetDuration("timeout"),
		version:                 version,
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

func (c *Config) Version() string {
	return c.version
}

func (c *Config) WebListenAddress() string {
	return c.webListenAddress
}
