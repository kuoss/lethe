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

	err := v.ReadInConfig()
	if err != nil {
		return &Config{}, fmt.Errorf("readInConfig err: %w", err)
	}
	err = viper.Unmarshal(&v)
	if err != nil {
		return &Config{}, fmt.Errorf("unmarshal err: %w", err)
	}

	logDataPath := v.GetString("storage.log_data_path")
	if logDataPath == "" {
		logDataPath = "./tmp/log"
	}

	retentionSize, err := util.StringToBytes(v.GetString("retention.size"))
	if err != nil {
		return &Config{}, fmt.Errorf("stringToBytes err: %w", err)
	}
	retentionTimeString := v.GetString("retention.time")
	if retentionTimeString == "" {
		retentionTimeString = "15d"
	}
	retentionTime, err := util.GetDurationFromAge(retentionTimeString)
	if err != nil {
		return &Config{}, fmt.Errorf("getDurationFromAge err: %w", err)
	}

	retentionSizingStrategy := v.GetString("retention.sizingStrategy")
	if retentionSizingStrategy == "" {
		retentionSizingStrategy = "file"
	}

	timeout := 20 * time.Second

	webListenAddress := v.GetString("web.listen_address")
	if webListenAddress == "" {
		webListenAddress = ":6060"
	}

	cfg := Config{
		limit:                   1000,
		logDataPath:             logDataPath,
		retentionSize:           retentionSize,
		retentionTime:           retentionTime,
		retentionSizingStrategy: retentionSizingStrategy,
		timeout:                 timeout,
		version:                 version,
		webListenAddress:        webListenAddress,
	}

	logger.Infof("====================================")
	logger.Infof("%+v", cfg)
	logger.Infof("====================================")
	return &cfg, nil
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
