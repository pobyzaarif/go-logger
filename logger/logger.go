package logger

import (
	"fmt"
	"runtime"
	"time"

	"github.com/labstack/gommon/log"
	goLoggerAppName "github.com/pobyzaarif/go-logger/appname"
)

var (
	logger = new()
	app    = goLoggerAppName.GetAPPName()
)

func new() *log.Logger {
	l := log.New("")
	l.DisableColor()
	l.SetHeader(`{"time":"${time_rfc3339_nano}","level":"${level}"}`)
	return l
}

type newLog struct {
	tag        string
	trackerID  string
	Caller     string // for manipulate or customizing caller value
	timerStart time.Time
}

func NewLog(tag string) newLog {
	return newLog{
		tag: tag,
	}
}

func (newLog *newLog) newLogParams(message string, data map[string]interface{}, err error) map[string]interface{} {
	logParams := make(map[string]interface{})

	if newLog.Caller == "" {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			logParams["caller"] = fmt.Sprintf("%v:%v", file, line)
		} else {
			logParams["caller"] = ""
		}
	} else {
		logParams["caller"] = newLog.Caller
	}

	logParams["service_name"] = app
	logParams["message"] = message
	logParams["tag"] = newLog.tag
	logParams["tracker_id"] = newLog.trackerID

	timeStart := time.Now()
	if !newLog.timerStart.IsZero() {
		timeStart = newLog.timerStart
	}

	elapsed := time.Since(timeStart)
	logParams["timer_start"] = timeStart
	logParams["timer_end"] = time.Now()
	logParams["processing_time"] = float64(elapsed.Nanoseconds() / 1e6)
	newLog.timerStart = time.Time{}

	if data != nil {
		// detect which one is gologger default
		if def, _ := data["__gologger__"].(int); def > 0 {
			delete(data, "__gologger__")
			logParams["data"] = data
		} else {
			logParams["data"] = map[string]interface{}{
				"app": data,
			}
		}
	} else {
		logParams["data"] = make(map[string]interface{})
	}

	if err != nil {
		logParams["error"] = err.Error()
	} else {
		logParams["error"] = ""
	}

	return logParams
}

func (newLog *newLog) TimerStart() {
	newLog.timerStart = time.Now()
}

func (newLog *newLog) SetTimerStart(timeStart time.Time) {
	newLog.timerStart = timeStart
}

func (newLog *newLog) SetTrackerID(trackerID string) {
	newLog.trackerID = trackerID
}

func (newLog *newLog) SetCallerValue(caller string) {
	newLog.Caller = caller
}

func (newLog *newLog) Info(message string) {
	logParams := newLog.newLogParams(message, nil, nil)
	logger.Infoj(logParams)
}

func (newLog *newLog) InfoWithData(message string, data map[string]interface{}) {
	logParams := newLog.newLogParams(message, data, nil)
	logger.Infoj(logParams)
}

func (newLog *newLog) Warn(message string) {
	logParams := newLog.newLogParams(message, nil, nil)
	logger.Warnj(logParams)
}

func (newLog *newLog) WarnWithData(message string, data map[string]interface{}) {
	logParams := newLog.newLogParams(message, data, nil)
	logger.Warnj(logParams)
}

func (newLog *newLog) WarnWithDataAndError(message string, data map[string]interface{}, err error) {
	logParams := newLog.newLogParams(message, data, err)
	logger.Warnj(logParams)
}

func (newLog *newLog) Error(message string, err error) {
	logParams := newLog.newLogParams(message, nil, err)
	logger.Errorj(logParams)
}

func (newLog *newLog) ErrorWithData(message string, data map[string]interface{}, err error) {
	logParams := newLog.newLogParams(message, data, err)
	logger.Errorj(logParams)
}

func (newLog *newLog) Fatal(message string) {
	logParams := newLog.newLogParams(message, nil, nil)
	logger.Fatalj(logParams)
}

func (newLog *newLog) FatalWithData(message string, data map[string]interface{}) {
	logParams := newLog.newLogParams(message, data, nil)
	logger.Fatalj(logParams)
}

func (newLog *newLog) FatalWithDataAndError(message string, data map[string]interface{}, err error) {
	logParams := newLog.newLogParams(message, data, err)
	logger.Fatalj(logParams)
}
