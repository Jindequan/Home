package job

import (
	"github.com/robfig/cron/v3"
)


type CronJob struct {
}

func Load() (job *CronJob, err error) {
	cronList := cron.New()
	//_, err = cronList.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	//_, err = cronList.AddFunc("@hourly", func() { fmt.Println("Every hour, starting an hour from now") })
	cronList.Start()
	return
}


