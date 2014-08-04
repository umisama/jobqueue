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

	ErrQueueNotFound  error = errors.New("queue not found.")
	ErrJobQueueIsFull error = errors.New("this job queue is full.")
)

// QueueConfig reprecents configulation for a job queue.
// MsgContener and Name must unique type in all queues.
type QueueConfig struct {
	Name        string
	Func        interface{}
	MsgContener interface{}
	Concurrency int
	Length      int
}

type queueConfig struct {
	Concurrency int
	Func        reflect.Value
	MsgContener reflect.Type
	Ch          chan reflect.Value
}

// SetConfig sets new configulation for job queue.
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
	if c.Length <= 0 {
		c.Length = 100 // default
	}

	configList[c.Name] = queueConfig{
		Concurrency: c.Concurrency,
		Func:        reflect.ValueOf(c.Func),
		MsgContener: reflect.ValueOf(c.MsgContener).Type(),
		Ch:          make(chan reflect.Value, c.Length),
	}

	for i:=0; i<c.Concurrency; i++ {
		go listenAndInvoke(c.Name)
	}
}

// SetConfigList sets new configulations for job queue with c.
func SetConfigList(c []QueueConfig) error {
	for _, v := range c {
		err := SetConfig(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func Publish(job interface{}) error {
	jobval := reflect.ValueOf(job)

	for _, conf := range configList {
		if conf.MsgContener == jobval.Type() {
			return publish(conf, jobval)
		}
	}

	return ErrQueueNotFound
}

func publish(conf queueConfig, job reflect.Value) error {
	if len(conf.Ch) == cap(conf.Ch) {
		return ErrJobQueueIsFull
	}

	conf.Ch <- job
	return nil
}

func init() {
	configList = make(map[string]queueConfig)
}
