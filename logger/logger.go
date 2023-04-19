package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = logrus.New()
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		CallerPrettyfier: getCallerPrettyfier(),
	})
}

func getCallerPrettyfier() func(f *runtime.Frame) (string, string) {
	return func(f *runtime.Frame) (string, string) {
		// https://github.com/sirupsen/logrus/blob/v1.9.0/example_custom_caller_test.go
		// https://github.com/kubernetes/klog/blob/v2.90.1/klog.go#L644
		_, file, line, ok := runtime.Caller(10)
		if !ok {
			file = "???"
			line = 1
		} else {
			if slash := strings.LastIndex(file, "/"); slash >= 0 {
				file = file[slash+1:]
			}
		}
		return "", fmt.Sprintf("%s:%d", file, line)
	}
}

func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
