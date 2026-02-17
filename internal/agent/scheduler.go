package agent

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var cronParser = cron.NewParser(
	cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

func shouldRun(job Job, now time.Time) bool {
	if !job.Enabled {
		return false
	}

	if job.LockUntil != nil && job.LockUntil.After(now) {
		return false
	}

	schedule, err := cronParser.Parse(job.Cron)
	if err != nil {
		log.Printf("invalid cron for job %s: %v", job.ID, err)
		return false
	}

	var last time.Time
	if job.LastRunAt != nil {
		last = *job.LastRunAt
	} else {
		last = now.Add(-24 * time.Hour)
	}

	next := schedule.Next(last)

	return now.After(next) || now.Equal(next)
}
