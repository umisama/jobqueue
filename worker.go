package jq

import (
	"reflect"
)

func listenAndInvoke(conf queueConfig) {
	for jobcnt := range conf.Ch {
		if len(conf.KillCh) != 0 { // this queue was killed.
			<-conf.KillCh
			return
		}

		jobcnt.info.Status = StatusRunning

		// call job
		rets := conf.Func.Call([]reflect.Value{jobcnt.job})
		jobcnt.info.Status = StatusCompleted
		if len(rets) != 0 {
			jobcnt.info.Return = rets[0].Interface()
		}
	}
}
