package jq

import (
	"errors"
	"reflect"
)

var (
	configList map[string]queueConfig

	ErrQueueIsExist        error = errors.New("this queue is exist.")
	ErrContenerIsNotUnique error = errors.New("c.Contener is not unique.")

	ErrQueueNotFound  error = errors.New("queue not found.")
	ErrJobQueueIsFull error = errors.New("this job queue is full.")
)

// QueueConfig reprecents configulation for a job queue.
// JobContener and Name must unique type in all queues.
type QueueConfig struct {
	Name        string
	JobContener Job
	Concurrency int
	Length      int
}

type queueConfig struct {
	Concurrency int
	Contener    reflect.Type
	Ch          chan jobContener
	KillCh      chan struct{}
}

// JobStatus reprecents status of a job.
type JobStatus int

const (
	StatusWaiting = JobStatus(iota)
	StatusRunning
	StatusCompleted
	StatusFailed
)

// JobInfomation is the infomation of job.
// Return is job's return value if status is completed. Error is in failed.
type JobInfomation struct {
	Id     string
	Status JobStatus
	Result interface{}
	Error  error
	Done   chan struct{}
}

type jobContener struct {
	job  Job
	info *JobInfomation
}

type Job interface {
	Run() error
	Result() interface{}
}

// SetConfig sets new configulation for job queue.
func SetConfig(c QueueConfig) error {
	// is name unique?
	if _, ok := configList[c.Name]; ok {
		return ErrQueueIsExist
	}

	// is contener unique?
	for _, v := range configList {
		if v.Contener == reflect.ValueOf(c.JobContener).Type() {
			return ErrContenerIsNotUnique
		}
	}

	return setConfig(c)
}

func setConfig(c QueueConfig) error {
	if c.Length <= 0 {
		c.Length = 100 // default
	}
	if c.Concurrency <= 0 {
		c.Concurrency = 3 // default
	}

	conf := queueConfig{
		Concurrency: c.Concurrency,
		Contener:    reflect.ValueOf(c.JobContener).Type(),
		Ch:          make(chan jobContener, c.Length),
		KillCh:      make(chan struct{}, 1),
	}
	configList[c.Name] = conf

	// run worker
	for i := 0; i < c.Concurrency; i++ {
		go listenAndInvoke(conf)
	}

	return nil
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

func Publish(job Job) (*JobInfomation, error) {
	for _, conf := range configList {
		if conf.Contener == reflect.ValueOf(job).Type() {
			return publish(conf, job)
		}
	}

	return nil, ErrQueueNotFound
}

func publish(conf queueConfig, job Job) (*JobInfomation, error) {
	if len(conf.Ch) == cap(conf.Ch) {
		return nil, ErrJobQueueIsFull
	}

	info := &JobInfomation{
		Id:     uuid(),
		Status: StatusWaiting,
		Result: nil,
		Error:  nil,
		Done:   make(chan struct{}, 0),
	}

	conf.Ch <- jobContener{
		job:  job,
		info: info,
	}
	return info, nil
}

func Kill(name string) error {
	c, ok := configList[name]
	if !ok {
		return ErrQueueNotFound
	}

	c.KillCh <- struct{}{}
	delete(configList, name)
	return nil
}

func init() {
	configList = make(map[string]queueConfig)
}
