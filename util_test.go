package jq

import (
	"testing"
)

func TestUuid(t *testing.T) {
	uuid1 := uuid()
	uuid2 := uuid()
	t.Log(uuid1, uuid2)
	if uuid1 == uuid2 {
		t.Error("")
	}

	if len(uuid1) != 36 || len(uuid2) != 36 {
		t.Error("")
	}
}
