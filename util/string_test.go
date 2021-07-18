package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStrExistsInFalse(t *testing.T) {
	needle := "1"
	haystackFalse := []string{"10", "20", "30"}
	expected := false
	found := StrExistsIn(needle, haystackFalse)
	if expected != found {
		t.Error(fmt.Sprintf("did not expect to find needle '%s' found in '%s'", needle, haystackFalse))
	}
}

func TestStrExistsInTrue(t *testing.T) {
	needle := "1"
	haystackFalse := []string{"10", "20", "30", "1"}
	expected := true
	found := StrExistsIn(needle, haystackFalse)
	if expected != found {
		t.Error(fmt.Sprintf("did expect to find needle '%s' found in '%s'", needle, haystackFalse))
	}
}

func TestStrOrSecond(t *testing.T) {
	first, second := "", "something"
	expected := second
	result := StrOr(first, second)
	if expected != result {
		t.Error("expected to find", expected, "and not", result)
	}
}

func TestStrOrFirst(t *testing.T) {
	first, second := "something", "something-else"
	expected := first
	result := StrOr(first, second)
	if expected != result {
		t.Error("expected to find", expected, "and not", result)
	}
}

func TestStrArrToChEmpty(t *testing.T) {
	in := []string{}
	expected := []string{}
	c := StrArrToCh(in)
	result := make([]string, 0)
	for next := range c {
		result = append(result, next)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Error("expected", expected, "and got", result)
	}
}

func TestStrArrToChFull(t *testing.T) {
	in := []string{"1", "2", "3"}
	expected := []string{"1", "2", "3"}
	c := StrArrToCh(in)
	result := make([]string, 0)
	for next := range c {
		result = append(result, next)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Error("expected", expected, "and got", result)
	}
}

func TestStrArrToChWithEmpty(t *testing.T) {
	in := []string{"1", "", "3"}
	expected := []string{"1", "3"}
	c := StrArrToCh(in)
	result := make([]string, 0)
	for next := range c {
		result = append(result, next)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Error("expected", expected, "and got", result)
	}
}
