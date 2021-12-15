package middleware

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	goutilAppName "github.com/pobyzaarif/goutil/appname"
	goutilHttp "github.com/pobyzaarif/goutil/http"
	goutilLogger "github.com/pobyzaarif/goutil/logger"
)

var (
	app                = goutilAppName.GetAPPName()
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
	reqHeader := goutilHttp.DumpRequest(c.Request(), []string{"Authorization"})

	logger := goutilLogger.NewLog("inbound_request")
	logger.SetTimerStart(reqTime)
	logger.SetTrackerID(fmt.Sprintf("%v", c.Get("tracker_id")))
	logger.InfoWithData("call_in", map[string]interface{}{
		"package":            packHandler,
		"handler":            funcHandler,
		"remote_ip":          c.RealIP(),
		"host":               c.Request().Host,
		"method":             c.Request().Method,
		"url":                c.Request().RequestURI,
		"request_time":       c.Request().Header.Get(headerRequestTime),
		"request_header":     reqHeader,
		"request":            string(req),
		"response_time":      c.Response().Header().Get(headerResponseTime),
		"response_header":    string(respHeader),
		"response":           string(res),
		"response_http_code": c.Response().Status,
	})
}
