package service

import (
	"log"

	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

func StartScheduler() {
	cronScheduler = cron.New()

	cronExpr := getAccessKeyCron()
	if cronExpr != "" {
		_, err := cronScheduler.AddFunc(cronExpr, DisableExpiredAccessKeys)
		if err != nil {
			log.Printf("scheduler: invalid cron expression '%s': %v", cronExpr, err)
		} else {
			log.Printf("scheduler: access_key expiry check scheduled with cron '%s'", cronExpr)
		}
	}

	cronScheduler.Start()
}

func ReloadAccessKeyCron() {
	if cronScheduler == nil {
		return
	}

	for _, entry := range cronScheduler.Entries() {
		cronScheduler.Remove(entry.ID)
	}

	cronExpr := getAccessKeyCron()
	if cronExpr != "" {
		_, err := cronScheduler.AddFunc(cronExpr, DisableExpiredAccessKeys)
		if err != nil {
			log.Printf("scheduler: invalid cron expression '%s': %v", cronExpr, err)
		} else {
			log.Printf("scheduler: access_key expiry check reloaded with cron '%s'", cronExpr)
		}
	}
}
