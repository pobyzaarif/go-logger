package logger

import (
	"context"
	"fmt"

	goLoggerDB "github.com/pobyzaarif/go-logger/database"
	goLogger "github.com/pobyzaarif/go-logger/logger"
	"go.mongodb.org/mongo-driver/event"
)

func Monitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			ctxTrackerID := ctx.Value("tracker_id")
			trackerID := ""
			if ctxTrackerID != nil {
				trackerID = fmt.Sprintf("%v", ctxTrackerID)
			}

			logger := goLogger.NewLog("MONGO_QUERY")
			logger.SetTrackerID(trackerID)

			logger.InfoWithData("query_info", goLoggerDB.DatabaseLog(map[string]interface{}{
				"query": evt.Command.String(),
			}))
		},
	}
}
