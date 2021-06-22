package main

import (
	"fmt"
	"testing"
)

func silentOut(format string, msg ...interface{}) {
}

func spyIntoString(target *string) outFunc {
	return func(format string, msg ...interface{}) {
		*target += fmt.Sprintf(format, msg...)
	}
}

func TestNow(t *testing.T) {
	var out, errOut string
	status := run(parameters{now: true, flagsSet: 1}, spyIntoString(&out), spyIntoString(&errOut))
	if status != 0 {
		t.Error("unexpected status")
	}
	if errOut != "" {
		t.Error("should not have panicked")
	}
	if out == "" {
		t.Error("should have created something")
	}
}
