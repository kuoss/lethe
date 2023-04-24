package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kuoss/common/logger"
	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var config *viper.Viper
var writer io.Writer
var limit = 1000
var logRoot = "./tmp/log"

func LoadConfig() error {
	config = viper.New()
	config.SetConfigName("lethe")
	config.SetConfigType("yaml")
	config.AddConfigPath(filepath.Join(".", "etc"))
	config.AddConfigPath(filepath.Join("..", "etc"))
	err := config.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error on ReadInConfig: %w", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("error on Unmarshal: %w", err)
	}
	SetLogRoot(config.GetString("storage.rootdirectory"))

	// show all settings in yaml format
	yamlBytes, err := yaml.Marshal(config.AllSettings())
	if err != nil {
		return fmt.Errorf("error on Marshal: %w", err)
	}
	logger.Infof("settings:\n====================================\n" + string(yamlBytes) + "====================================")
	return nil
}

func GetConfig() *viper.Viper {
	return config
}

// FOR CLI
func SetWriter(w io.Writer) {
	writer = w
}

func GetWriter() io.Writer {
	if writer == nil {
		return os.Stdout
	}
	return writer
}

func GetLimit() int {
	return limit
}

func SetLimit(newLimit int) {
	limit = newLimit
}

func GetLogRoot() string {
	return logRoot
}

func SetLogRoot(newLogRoot string) {
	logRoot = newLogRoot
}
