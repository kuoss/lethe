package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "github.com/kuoss/lethe/storage/driver/filesystem"
	"github.com/spf13/viper"
)

var config *viper.Viper
var writer io.Writer
var limit = 1000
var logRoot = "./tmp/log"

func LoadConfig() {
	config = viper.New()
	config.SetConfigName("lethe")
	config.SetConfigType("yaml")
	config.AddConfigPath(filepath.Join(".", "etc"))
	config.AddConfigPath(filepath.Join("..", "etc"))
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config logs: %w", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to decode into struct, %v", err))
	}

	SetLogRoot(config.GetString("storage.rootdirectory"))
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
