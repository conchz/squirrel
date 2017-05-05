package log

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/lavenderx/squirrel/app"
	"go.uber.org/zap"
	"strings"
	//"sync"
	"github.com/lavenderx/squirrel/app/log/config"
	"gopkg.in/yaml.v2"
)

// https://github.com/emonuh/zap-error-dispatcher

var (
	logger = new(zap.Logger)
	//logger = new(zap.SugaredLogger)
	//lock   = new(sync.RWMutex)
)

func GetLogger() *zap.Logger {
	return logger
}

//func Debug(args ...interface{}) {
//	logger.Debug(args...)
//}
//
//func Debugf(template string, args ...interface{}) {
//	logger.Debugf(template, args)
//}
//
//func Info(args ...interface{}) {
//	logger.Info(args...)
//}
//
//func Infof(template string, args ...interface{}) {
//	logger.Infof(template, args)
//}
//
//func Warn(args ...interface{}) {
//	logger.Warn(args...)
//}
//
//func Warnf(template string, args ...interface{}) {
//	logger.Warnf(template, args)
//}
//
//func Error(args ...interface{}) {
//	logger.Error(args...)
//}
//
//func Errorf(template string, args ...interface{}) {
//	logger.Errorf(template, args)
//}
//
//func Fatal(args ...interface{}) {
//	logger.Fatal(args...)
//}
//
//func Fatalf(template string, args ...interface{}) {
//	logger.Fatalf(template, args)
//}
//
//func Panic(args ...interface{}) {
//	logger.Panic(args...)
//}
//
//func Panicf(template string, args ...interface{}) {
//	logger.Panicf(template, args)
//}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) log.Lvl {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	case "OFF":
		return log.OFF
	}

	panic(fmt.Errorf("not a valid log Level: %q", lvl))
}

func Init() {
	var _zapConfig config.ErrorDispatcherConfig
	if err := yaml.Unmarshal(app.GetLogConfBytes(), &_zapConfig); err != nil {
		panic(err)
	}

	var err error
	logger, err = _zapConfig.Build()
	if err != nil {
		panic(err)
	}

	//logConf := app.LoadConfig().LoggingConf
	//
	//zapConfig := zap.NewProductionConfig()
	//zapConfig.EncoderConfig.EncodeTime.UnmarshalText([]byte("ISO8601"))
	//zapConfig.DisableStacktrace = true
	//zapConfig.Level.UnmarshalText([]byte(logConf.Level))
	//zapLogger, err := zapConfig.Build()
	//if err != nil {
	//	fmt.Errorf("%v", err)
	//	return
	//}
	//
	//lock.Lock()
	//defer lock.Unlock()
	//logger = zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}
