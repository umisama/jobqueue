package jq

import (
	"reflect"
)

func listenAndInvoke(name string) {
	conf := configList[name]
	for job := range conf.Ch {
		conf.Func.Call([]reflect.Value{job})
	}
}
