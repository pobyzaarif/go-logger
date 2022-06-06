package middleware

import (
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	goLoggerAppName "github.com/pobyzaarif/go-logger/appname"
	goLoggerHttp "github.com/pobyzaarif/go-logger/http"
	goLogger "github.com/pobyzaarif/go-logger/logger"
)

var (
	app                = goLoggerAppName.GetAPPName()
	headerRequestTime  = "X-" + app + "-RequestTime"
	headerResponseTime = "X-" + app + "-ResponseTime"
)

func ServiceRequestTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set(headerRequestTime, time.Now().Format(time.RFC3339Nano))
		return next(c)
	}
}

func ServiceTrackerID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("tracker_id", uuid.New().String())
		return next(c)
	}
}

func APILogHandler(c echo.Context, req, res []byte) {
	c.Response().Header().Set(headerResponseTime, time.Now().Format(time.RFC3339Nano))
	reqTime, err := time.Parse(time.RFC3339Nano, c.Request().Header.Get(headerRequestTime))
	if err != nil {
		reqTime = time.Now()
	}

	var handler string
	r := c.Echo().Routes()
	cpath := strings.Replace(c.Path(), "/", "", -1)
	for _, v := range r {
		vpath := strings.Replace(v.Path, "/", "", -1)
		if vpath == cpath && v.Method == c.Request().Method {
			handler = v.Name
			// Handler for wrong route.
			if strings.Contains(handler, "func1") {
				handler = "UndefinedRoute"
			}
			break
		}
	}

	// Get Handler Name
	dir, file := path.Split(handler)
	fileStrings := strings.Split(file, ".")
	packHandler := dir + fileStrings[0]
	funcHandler := strings.Replace(handler, packHandler+".", "", -1)

	respHeader, _ := json.Marshal(c.Response().Header())
	reqHeader := goLoggerHttp.DumpRequest(c.Request(), []string{"Authorization"})

	tranckerID, _ := c.Get("tracker_id").(string)
	logger := goLogger.NewLog("INBOUND_REQUEST")
	logger.SetTimerStart(reqTime)
	logger.SetTrackerID(tranckerID)
	logger.InfoWithData("api_info", goLoggerHttp.NetworkLog(map[string]interface{}{
		"handler":            funcHandler,
		"remote_ip":          c.RealIP(),
		"host":               c.Request().Host,
		"method":             c.Request().Method,
		"url":                c.Request().RequestURI,
		"request_header":     reqHeader,
		"request":            string(req),
		"response_header":    string(respHeader),
		"response":           string(res),
		"response_http_code": c.Response().Status,
	}))
}
