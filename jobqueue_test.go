package jq

import (
	"testing"
	"time"
)

type TestJob1 struct {
	Ch chan struct{}
}

func (j *TestJob1) Run() {
	j.Ch <- struct{}{}
}

func (j *TestJob1) Result() interface{}{
	return nil
}

type TestJob2 struct {
	Ch chan struct{}
}

func (j *TestJob2) Run() {
	j.Ch <- struct{}{}
}

func (j *TestJob2) Result() interface{} {
	return nil
}

func TestSetConfig(t *testing.T) {
	type testcase struct {
		in     QueueConfig
		expect error
	}
	cases := []testcase{{
		in: QueueConfig{
			Name:        "testname",
			Concurrency: 1,
			JobContener: &TestJob1{},
		},
		expect: nil,
	}, {
		in: QueueConfig{
			Name:        "testname",
			Concurrency: 1,
			JobContener: &TestJob1{},
		},
		expect: ErrQueueIsExist,
	}, {
		in: QueueConfig{
			Name:        "testname2",
			Concurrency: 1,
			JobContener: &TestJob1{},
		},
		expect: ErrContenerIsNotUnique,
	}}

	for _, v := range cases {
		err := SetConfig(v.in)
		if err != v.expect {
			t.Errorf("expect %s but got %s", v.expect, err)
		}
	}

	// clear
	clearutil()
}

func TestSetConfigList(t *testing.T) {
	type testcase struct {
		in     []QueueConfig
		expect error
	}
	cases := []testcase{{
		in: []QueueConfig{{
			Name:        "testname",
			Concurrency: 1,
			JobContener: &TestJob1{},
		}, {
			Name:        "testname2",
			Concurrency: 1,
			JobContener: &TestJob2{},
		}},
		expect: nil,
	}}

	for _, v := range cases {
		err := SetConfigList(v.in)
		if err != v.expect {
			t.Errorf("expect %s but got %s", v.expect, err)
		}
	}
	// clear
	clearutil()
}

func TestPublish(t *testing.T) {
	// initialize
	SetConfig(QueueConfig{
		Name:        "test1",
		JobContener: &TestJob1{},
		Concurrency: 1,
		Length:      1,
	})

	// returns error if job is not registerd.
	_, err := Publish(&TestJob2{})
	if err != ErrQueueNotFound {
		t.Fail()
	}

	// returns info if job is registerd.
	ch := make(chan struct{}, 0)
	info, err := Publish(&TestJob1{ch})
	if err != nil {
		t.Fail()
	}
	if info == nil {
		t.Fail()
	}

	// returns error if queue is full
	_, err = Publish(&TestJob1{ch})
	if err != ErrJobQueueIsFull {
		t.Fail()
	}

	// must run(1st called job)
	select {
	case <-time.After(10 * time.Second):
		t.Fail()
	case <-ch:
		t.Log("ok")
	}

	Kill("test1")
}

func clearutil() {
	configList = map[string]queueConfig{}
}
