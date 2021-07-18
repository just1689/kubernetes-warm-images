package util

import (
	"fmt"
	"testing"
)

func TestLogPrependL1(t *testing.T) {
	msg := "xyz"
	expectedLevel1 := " >   : " + msg
	result := LogPrepend(1, msg)
	if expectedLevel1 != result {
		t.Error(fmt.Sprintf("expected to find '%s' and not '%s'", expectedLevel1, result))
	}
}

func TestLogPrependL2(t *testing.T) {
	msg := "xyz"
	expectedLevel1 := " >>  : " + msg
	result := LogPrepend(2, msg)
	if expectedLevel1 != result {
		t.Error(fmt.Sprintf("expected to find '%s' and not '%s'", expectedLevel1, result))
	}
}

func TestLogPrependL3(t *testing.T) {
	msg := "xyz"
	expectedLevel1 := "     : " + msg
	result := LogPrepend(3, msg)
	if expectedLevel1 != result {
		t.Error(fmt.Sprintf("expected to find '%s' and not '%s'", expectedLevel1, result))
	}
}
