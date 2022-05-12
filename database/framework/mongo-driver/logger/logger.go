package logger

import (
	"context"
	"fmt"

	goutilLogger "github.com/pobyzaarif/goutil/logger"
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

			logger := goutilLogger.NewLog("MONGO_COMMAND")
			logger.SetTrackerID(trackerID)

			logger.InfoWithData("command_info", map[string]interface{}{
				"command": evt.Command.String(),
			})
		},
	}
}
