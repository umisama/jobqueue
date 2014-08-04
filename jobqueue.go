package jq

import (
	"errors"
	"reflect"
)

var (
	configList map[string]queueConfig

	ErrQueueIsExist        error = errors.New("this queue is exist.")
	ErrFuncIsNotFunction   error = errors.New("c.Func is not function.")
	ErrContenerIsNotUnique error = errors.New("c.Contener is not unique.")

	ErrQueueNotFound error = errors.New("error queue not found.")
)

// QueueConfig reprecents configulation for a job queue.
// MsgContener must unique type in all queues.
type QueueConfig struct {
	Name        string
	Func        interface{}
	MsgContener interface{}
	Concurrency int
	Length      int
}

type queueConfig struct {
	Concurrency int
	Func        interface{}
	MsgContener reflect.Type
	Ch          chan interface{}
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
		Ch:          make(chan interface{}, c.Length),
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
	configList = make(map[string]queueConfig)
}
