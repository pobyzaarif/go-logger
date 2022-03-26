package middleware

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	goutilLogger "github.com/pobyzaarif/goutil/logger"
)

type (
	Skipper func(echo.Context) bool
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `yaml:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `yaml:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `yaml:"disable_print_stack"`

		// LogLevel is log level to printing stack trace.
		// Optional. Default value 0 (Print).
		// LogLevel log.Lvl
	}
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(echo.Context) bool {
	return false
}

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	DefaultRecoverConfig = RecoverConfig{
		Skipper:           DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		// LogLevel:          0,
	}
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	return recoverWithConfig(DefaultRecoverConfig)
}

// recoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func recoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						tranckerID, _ := c.Get("tracker_id").(string)
						logger := goutilLogger.NewLog("PANIC")
						logger.SetTrackerID(tranckerID)

						msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])

						logger.ErrorWithData("PANIC", map[string]interface{}{
							"trace": msg,
						}, err)
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
