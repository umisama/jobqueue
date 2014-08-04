package jq

import (
	"testing"
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
		},{
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

func clearutil() {
	configList = map[string]queueConfig{}
}
