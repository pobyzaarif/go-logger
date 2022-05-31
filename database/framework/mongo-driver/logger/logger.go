package logger

import (
	"context"
	"fmt"

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

			logger := goLogger.NewLog("MONGO_COMMAND")
			logger.SetTrackerID(trackerID)

			logger.InfoWithData("command_info", map[string]interface{}{
				"command": evt.Command.String(),
			})
		},
	}
}
