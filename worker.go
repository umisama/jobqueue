package jq

import (
	"errors"
	"fmt"
)

func listenAndInvoke(conf queueConfig) {
	for cnt := range conf.Ch {
		if len(conf.KillCh) != 0 { // this queue was killed.
			<-conf.KillCh
			return
		}

		cnt.info.Status = StatusRunning
		err := invoke(cnt.job)
		if err != nil {
			cnt.info.Status = StatusFailed
			cnt.info.Error = err
			continue
		}

		cnt.info.Status = StatusCompleted
		cnt.info.Result = cnt.job.Result()
	}
}

func invoke(job Job)(err error){
	defer func() {
		if r := recover(); r != nil {
			switch rt := r.(type) {
			case error:
				err = rt
			case string:
				err = errors.New(rt)
			default:
				err = fmt.Errorf("%v", rt)
			}
		}
	}()

	err = job.Run()
	if err != nil {
		return err
	}
	return nil
}
