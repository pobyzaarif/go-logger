package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	goutilLogger "github.com/pobyzaarif/goutil/logger"
	lg "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var ErrRecordNotFound = errors.New("record not found")

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LogLevel
type LogLevel int

const (
	Silent LogLevel = iota + 1
	Error
	Warn
	Info
)

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

type Config struct {
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	LogLevel                  lg.LogLevel
}

// lg.Interface logger interface
type Interface interface {
	LogMode(lg.LogLevel) lg.Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

var (
	Default = New(log.New(os.Stdout, "\r\n", log.LstdFlags), Config{
		// SlowThreshold:             200 * time.Millisecond,
		SlowThreshold:             5 * time.Second,
		LogLevel:                  lg.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  false,
	})
)

func New(writer Writer, config Config) lg.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &logger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logger struct {
	Writer
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logger) LogMode(level lg.LogLevel) lg.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lg.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	ctxTrackerID := ctx.Value("tracker_id")
	trackerID := ""
	if ctxTrackerID != nil {
		trackerID = fmt.Sprintf("%v", ctxTrackerID)
	}

	logger := goutilLogger.NewLog("GORM_QUERY")
	logger.SetCallerValue(utils.FileWithLineNum())
	logger.SetTrackerID(trackerID)
	logger.SetTimerStart(begin)

	if l.LogLevel <= lg.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= lg.Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.ErrorWithData(
				"query_error",
				map[string]interface{}{
					"rows":  "-",
					"query": sql,
				},
				err,
			)
		} else {
			logger.ErrorWithData(
				"error_query",
				map[string]interface{}{
					"rows":  rows,
					"query": sql,
				},
				err,
			)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= lg.Warn:
		sql, rows := fc()
		if rows == -1 {
			logger.WarnWithData(
				"warn_slow_query",
				map[string]interface{}{
					"rows":  "-",
					"query": sql,
				},
			)
		} else {
			logger.WarnWithData(
				"warn_slow_query",
				map[string]interface{}{
					"rows":  rows,
					"query": sql,
				},
			)
		}
	case l.LogLevel == lg.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.InfoWithData(
				"query_info",
				map[string]interface{}{
					"rows":  "-",
					"query": sql,
				},
			)
		} else {
			logger.InfoWithData(
				"query_info",
				map[string]interface{}{
					"rows":  rows,
					"query": sql,
				},
			)
		}
	}
}

type traceRecorder struct {
	lg.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}
