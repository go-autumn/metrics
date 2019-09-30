package metrics

import "gopkg.in/robfig/cron.v2"

var defaultCron *cron.Cron

func init() {
	defaultCron = cron.New()
	defaultCron.Start()
}

func DefaultCron() *cron.Cron {
	return defaultCron
}
