package app

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
	"runtime"
)

// Extend echo's middleware

type (
	CustomContext struct {
		echo.Context
		logger *zap.SugaredLogger
	}
)

func (c *CustomContext) Logger() *zap.SugaredLogger {
	return c.logger
}

// Customize echo's RecoverWithConfig to use zap log
func RecoverWithConfig(config middleware.RecoverConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = middleware.DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						c.Get("logger").(*zap.SugaredLogger).Errorf("[%s] %s %s\n", "PANIC RECOVER", err, stack[:length])
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
