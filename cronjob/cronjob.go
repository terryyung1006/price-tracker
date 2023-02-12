package cronjob

import "time"

type Icronjob interface {
	Run(interval time.Duration)
}

func RunCronJob(cronJob Icronjob, interval time.Duration) {
	go func() {
		for {
			cronJob.Run(interval)
		}
	}()
}
