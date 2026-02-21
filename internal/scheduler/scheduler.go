package scheduler

import (
	"time"

	"DataFPP_Dingtalk/internal/service"
	"DataFPP_Dingtalk/pkg/logger"

	"github.com/go-co-op/gocron/v2"
)

func Setup(syncService *service.SyncService, timezone string) (gocron.Scheduler, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	s, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		return nil, err
	}

	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(8, 0, 0))),
		gocron.NewTask(func() {
			if err := syncService.SyncData(); err != nil {
				logger.Error("Scheduled sync failed: %v", err)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}
