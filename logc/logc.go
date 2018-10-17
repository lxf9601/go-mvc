// 日志封装
package logc

import (
	"os"
	"github.com/gogap/logrus"
)

func InitLogger(appPath string, logLevel string) {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: "01-02 15:04"})
	logrus.SetOutput(os.Stdout)
	logrus.AddHook(NewHook(appPath + "logs/ss.log"))
	if logLevel == "DEBUG" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if logLevel == "INFO" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if logLevel == "WARN" {
		logrus.SetLevel(logrus.WarnLevel)
	} else if logLevel == "ERROR" {
		logrus.SetLevel(logrus.ErrorLevel)
	} else if logLevel == "FATAL" {
		logrus.SetLevel(logrus.FatalLevel)
	} else if logLevel == "PANIC" {
		logrus.SetLevel(logrus.PanicLevel)
	}
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func IsDebug() bool {
	return logrus.GetLevel() >= logrus.DebugLevel
}

func IsInfo() bool {
	return logrus.GetLevel() >= logrus.InfoLevel
}
