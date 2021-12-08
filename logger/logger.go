package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000"
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

type eventLog struct {
	event             string
	defaultTimerStart time.Time
	timerStart        time.Time
}

func NewLog(event string) eventLog {
	return eventLog{
		event:             event,
		defaultTimerStart: time.Now(),
	}
}

func (eventLog *eventLog) Info(message string) {
	logParams := eventLog.newLogParams(map[string]interface{}{}, nil)
	logger.Info().Fields(logParams).Msg(message)
}

func (eventLog *eventLog) InfoWithData(message string, data map[string]interface{}) {
	logParams := eventLog.newLogParams(data, nil)
	logger.Info().Fields(logParams).Msg(message)
}

func (eventLog *eventLog) newLogParams(data map[string]interface{}, err error) map[string]interface{} {
	logParams := make(map[string]interface{})

	_, file, line, ok := runtime.Caller(2)
	if ok {
		data["caller"] = fmt.Sprintf("%v:%v", path.Base(file), line)
	}

	logParams["event"] = eventLog.event

	timeStart := eventLog.defaultTimerStart
	if !eventLog.timerStart.IsZero() {
		timeStart = eventLog.timerStart
	}

	logParams["timer_start"] = timeStart
	logParams["timer_end"] = time.Now()
	logParams["processing_time"] = time.Since(timeStart) / 1e3
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
