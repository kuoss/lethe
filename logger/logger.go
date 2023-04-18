package logger

import (
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	loggerInstance *logrus.Logger
	loggerOnce     sync.Once
)

// get singleton logger
func GetLogger() *logrus.Logger {
	if loggerInstance == nil {
		loggerOnce.Do(
			func() {
				loggerInstance = newLogger()
			})
	}
	return loggerInstance
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	})
	return logger
}
