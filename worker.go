package jq

import (
	"reflect"
)

func listenAndInvoke(conf queueConfig) {
	for job := range conf.Ch {
		if len(conf.KillCh) != 0 { // this queue was killed.
			<-conf.KillCh
			return
		}
		conf.Func.Call([]reflect.Value{job})
	}
}
