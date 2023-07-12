package logutil

import (
	"github.com/qiqiuyang/logger"
	"go.uber.org/zap"
)

var (
	logutil      logger.LoggerService
	suffixLogger = "DEFAULT"
	logName      = "default.log"
)

func init() {
	logutil = logger.NewLoggerService(nil)
	config := logutil.MakeDefaultLogConfig("", logName, suffixLogger)
	logutil.MakeLogger(config)
}

func GetSugarLogger() *zap.SugaredLogger {
	log, _ := logutil.GetSugarLogger(logName)
	return log
}

func GetLogger() *zap.Logger {
	log, _ := logutil.GetLogger(logName)
	return log
}
