package jq

import (
	"testing"
	"time"
)

func TestSetConfig(t *testing.T) {
	type testcase struct {
		in     QueueConfig
		expect error
	}
	testfunc := func() {}
	cases := []testcase{{
		in: QueueConfig{
			Name:        "testname",
			Concurrency: 1,
			Func:        testfunc,
			MsgContener: testcase{}, // as contener sample
		},
		expect: nil,
	}, {
		in: QueueConfig{
			Name:        "testnamehoge",
			Concurrency: 1,
			Func:        "this is not function",
			MsgContener: testcase{}, // as contener sample
		},
		expect: ErrFuncIsNotFunction,
	}, {
		in: QueueConfig{
			Name:        "testname",
			Concurrency: 1,
			Func:        testfunc,
			MsgContener: testcase{},
		},
		expect: ErrQueueIsExist,
	}, {
		in: QueueConfig{
			Name:        "testname2",
			Concurrency: 1,
			Func:        testfunc,
			MsgContener: testcase{},
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
	testfunc := func() {}
	cases := []testcase{{
		in: []QueueConfig{{
			Name:        "testname",
			Concurrency: 1,
			Func:        testfunc,
			MsgContener: testcase{}, // as contener sample
		}, {
			Name:        "testname2",
			Concurrency: 1,
			Func:        testfunc,
			MsgContener: struct{}{},
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
	// jobs for unit testing
	type TestJobTypeDummy struct{}
	type TestJobType struct{
		test chan struct{}
	}
	testfunc := func(j TestJobType) {
		time.Sleep(1*time.Second)
		j.test <- struct{}{}
	}

	// initialize
	SetConfig(QueueConfig{
		Name:        "test1",
		Func:        testfunc,
		MsgContener: TestJobType{},
		Concurrency: 1,
		Length:      1,
	})

	// returns error if job is not registerd.
	_, err := Publish(TestJobTypeDummy{})
	if err != ErrQueueNotFound {
		t.Fail()
	}

	// returns info if job is registerd.
	ch := make(chan struct{}, 0)
	info, err := Publish(TestJobType{ch})
	if err != nil {
		t.Fail()
	}
	if info == nil {
		t.Fail()
	}

	// returns error if queue is full
	Publish(TestJobType{ch})
	_, err = Publish(TestJobType{ch})
	if err != ErrJobQueueIsFull {
		t.Fail()
	}

	// must run(1st called job)
	select {
	case <- time.After(10 *time.Second):
		t.Fail()
	case <-ch:
		t.Log("ok")
	}

	Kill("test1")
}

func clearutil() {
	configList = map[string]queueConfig{}
}
