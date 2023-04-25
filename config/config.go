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

var (
	vip              *viper.Viper
	logDataPath      string
	webListenAddress string
	writer           io.Writer = os.Stdout
	limit            int       = 1000
)

func LoadConfig() error {
	vip = viper.New()

	vip.SetDefault("storage.log_data_path", "./tmp/log")
	vip.SetDefault("web.listen_address", ":6060")

	vip.SetConfigName("lethe")
	vip.SetConfigType("yaml")
	vip.AddConfigPath(filepath.Join(".", "etc"))
	vip.AddConfigPath(filepath.Join("..", "etc"))
	err := vip.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error on ReadInConfig: %w", err)
	}
	err = viper.Unmarshal(&vip)
	if err != nil {
		return fmt.Errorf("error on Unmarshal: %w", err)
	}
	SetLogDataPath(vip.GetString("storage.log_data_path"))
	SetWebListenAddress(vip.GetString("web.listen_address"))

	// show all settings in yaml format
	yamlBytes, err := yaml.Marshal(vip.AllSettings())
	if err != nil {
		return fmt.Errorf("error on Marshal: %w", err)
	}
	logger.Infof("settings:\n====================================\n" + string(yamlBytes) + "====================================")
	return nil
}

func Viper() *viper.Viper {
	return vip
}

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

func GetLogDataPath() string {
	return logDataPath
}

func SetLogDataPath(newLogDataPath string) {
	logDataPath = newLogDataPath
}

func GetWebListenAddress() string {
	return webListenAddress
}

func SetWebListenAddress(newWebListenAddress string) {
	webListenAddress = newWebListenAddress
}
