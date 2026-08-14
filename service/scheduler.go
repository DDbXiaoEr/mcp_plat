// Copyright (C) 2026 Zhaoquan Wang
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
