package jq

import (
	"errors"
	"reflect"
	"testing"
)

type PanicJob struct {
	val interface{}
}

func (j *PanicJob) Run() error {
	panic(j.val)
	return nil
}

func (j *PanicJob) Result() interface{} {
	return nil
}

func TestInvoke(t *testing.T) {
	type testcase struct {
		in     Job
		expect error
	}

	cases := []testcase{{
		in:     &PanicJob{errors.New("test")},
		expect: errors.New("test"),
	},{
		in:     &PanicJob{"string"},
		expect: errors.New("string"),
	},{
		in:     &PanicJob{123},
		expect: errors.New("123"),
	},{
		in:     &TestJob1{make(chan struct{}, 1)},
		expect: nil,
	}}

	for _, v := range cases {
		err := invoke(v.in)
		if !reflect.DeepEqual(err, v.expect) {
			t.Errorf("expect %s but got %s", v.expect, err)
		}
	}
}
