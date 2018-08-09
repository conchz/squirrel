package app

import (
	"github.com/labstack/echo"
	"github.com/lavenderx/squirrel/app/log"
	"github.com/lavenderx/squirrel/app/models"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"net/http"
	"strings"
)

func init() {
	// Setup app startup hooks
	OnAppStart(func() {
		// Load config at first
		config := LoadConfig()

		// Load assets and set assetsHandler
		assetsHandler := http.FileServer(Assets().HTTPBox())
		staticFilesHandler = echo.WrapHandler(assetsHandler)

		// Setup server port
		port = config.ServerConf.Port
	}, -1)
	OnAppStart(initLog, 0)
	OnAppStart(initMySQL)
	//OnAppStart(initRedis)

	// Setup app shutdown hooks
	OnAppStop(shutdownMySQL)
	//OnAppStop(shutdownRedis)
}

// Initialize log component
func initLog() {
	var zapLogConfig log.ZapLogConfig
	if err := yaml.Unmarshal(GetLogConfBytes(), &zapLogConfig); err != nil {
		panic(err)
	}
	logger = log.New(zapLogConfig)
	isDebug = strings.ToUpper(zapLogConfig.Level.Level().String()) == zapcore.DebugLevel.CapitalString()
}

// Initialize MySQL client
func initMySQL() {
	ConnectToMySQL(config)

	mySQLTemplate = GetMySQLTemplate()
	//if err := mySQLTemplate.XormEngine().Sync2(new(models.User)); err != nil {
	//	panic(err)
	//}
}

// Initialize Redis client
func initRedis() {
	ConnectToRedis(config)
}

// Close MySQL client
func shutdownMySQL() {
	logger.Info("Closing MySQL client")
	if err := GetMySQLTemplate().Close(); err != nil {
		logger.Error(err)
	}
}

// Close Redis client
func shutdownRedis() {
	logger.Info("Closing Redis client")
	if err := CloseRedisClient(); err != nil {
		logger.Error(err)
	}
}
