package app

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"strings"
)

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (log.Lvl, error) {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return log.DEBUG, nil
	case "INFO":
		return log.INFO, nil
	case "WARN":
		return log.WARN, nil
	case "ERROR":
		return log.ERROR, nil
	case "OFF":
		return log.OFF, nil
	}

	var l log.Lvl
	return l, fmt.Errorf("not a valid log Level: %q", lvl)
}
