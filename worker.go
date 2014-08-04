package jq

func listenAndInvoke(conf queueConfig) {
	for jobcnt := range conf.Ch {
		if len(conf.KillCh) != 0 { // this queue was killed.
			<-conf.KillCh
			return
		}

		jobcnt.info.Status = StatusRunning
		jobcnt.job.Run()
		jobcnt.info.Status = StatusCompleted
		jobcnt.info.Result = jobcnt.job.Result()
	}
}
