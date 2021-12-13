package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	goutilAppName "github.com/pobyzaarif/goutil/appname"
	"github.com/rs/zerolog"
)

var (
	logger      zerolog.Logger
	serviceName = "service_name"
	app         = goutilAppName.GetAPPName()
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000"
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

type eventLog struct {
	event             string
	trackerID         string
	defaultTimerStart time.Time
	timerStart        time.Time
}

func NewLog(event string) eventLog {
	return eventLog{
		event:             event,
		defaultTimerStart: time.Now(),
	}
}

func (eventLog *eventLog) newLogParams(data map[string]interface{}, err error) map[string]interface{} {
	logParams := make(map[string]interface{})

	_, file, line, ok := runtime.Caller(2)
	if ok {
		data["caller"] = fmt.Sprintf("%v:%v", path.Base(file), line)
	}

	logParams["event"] = eventLog.event
	logParams["tracker_id"] = eventLog.trackerID

	timeStart := eventLog.defaultTimerStart
	if !eventLog.timerStart.IsZero() {
		timeStart = eventLog.timerStart
	}

	logParams["timer_start"] = timeStart
	logParams["timer_end"] = time.Now()
	logParams["processing_time"] = time.Since(timeStart).Nanoseconds() / int64(time.Millisecond)
	eventLog.timerStart = time.Time{}

	if data != nil {
		logParams["data"] = data
	}

	if err != nil {
		logParams["error"] = err.Error()
	}

	return logParams
}

func (eventLog *eventLog) TimerStart() {
	eventLog.timerStart = time.Now()
}

func (eventLog *eventLog) SetTimerStart(timeStart time.Time) {
	eventLog.timerStart = timeStart
}

func (eventLog *eventLog) SetTrackerID(trackerID string) {
	eventLog.trackerID = trackerID
}

func (eventLog *eventLog) Info(message string) {
	logParams := eventLog.newLogParams(map[string]interface{}{}, nil)
	logger.Info().Str(serviceName, app).Fields(logParams).Msg(message)
}

func (eventLog *eventLog) InfoWithData(message string, data map[string]interface{}) {
	logParams := eventLog.newLogParams(data, nil)
	logger.Info().Str(serviceName, serviceName).Fields(logParams).Msg(message)
}
