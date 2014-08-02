package jq

import (
	"errors"
	"reflect"
)

var (
	queues     map[string]chan interface{}
	configList map[string]queueConfig

	ErrQueueIsExist        error = errors.New("this queue is exist.")
	ErrFuncIsNotFunction   error = errors.New("c.Func is not function.")
	ErrContenerIsNotUnique error = errors.New("c.Contener is not unique.")
)

// QueueConfig reprecents configulation for a job queue.
// MsgContener must unique type in all queues.
type QueueConfig struct {
	Name        string
	Concurrency int
	Func        interface{}
	MsgContener interface{}
}

type queueConfig struct {
	Concurrency int
	Func        interface{}
	MsgContener reflect.Type
}

// SetConfig() sets new QueueConfig
func SetConfig(c QueueConfig) error {
	// is name unique?
	if _, ok := configList[c.Name]; ok {
		return ErrQueueIsExist
	}

	// is c.Func function?
	if reflect.ValueOf(c.Func).Kind() != reflect.Func {
		return ErrFuncIsNotFunction
	}

	// is contener unique?
	cntval := reflect.ValueOf(c.MsgContener)
	for _, v := range configList {
		if v.MsgContener == cntval.Type() {
			return ErrContenerIsNotUnique
		}
	}

	setConfig(c)
	return nil
}

func setConfig(c QueueConfig) {
	configList[c.Name] = queueConfig{
		Concurrency: c.Concurrency,
		Func:        c.Func,
		MsgContener: reflect.ValueOf(c.MsgContener).Type(),
	}
}

func SetConfigList(c []QueueConfig) error {
	for _, v := range c {
		err := SetConfig(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	queues = make(map[string]chan interface{})
	configList = make(map[string]queueConfig)
}
