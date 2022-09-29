package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

// Config
// type Config struct {
// 	Retention Retention
// }
// type Retention struct {
// 	Time string
// 	Size string
// }

var config *viper.Viper
var writer io.Writer
var now time.Time
var limit = 1000
var logRoot = "/data/log"

func LoadConfig() {
	config = viper.New()
	config.SetConfigName("lethe")
	config.SetConfigType("yaml")
	config.AddConfigPath("./etc/")
	config.AddConfigPath("../etc/")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to decode into struct, %v", err))
	}
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

// FOR TEST
func SetNow(t time.Time) {
	now = t
}

func GetNow() time.Time {
	if reflect.ValueOf(now).IsZero() {
		return time.Now()
	}
	return now
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
