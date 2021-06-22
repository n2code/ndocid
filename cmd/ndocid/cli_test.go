package main

import (
	"fmt"
	"regexp"
	"testing"
)

func silentOut(format string, msg ...interface{}) {
}

func spyIntoString(target *string) outFunc {
	return func(format string, msg ...interface{}) {
		*target += fmt.Sprintf(format, msg...)
	}
}

func assertSuccess(p parameters, expRegex string, t *testing.T) {
	var out, errOut string
	status := run(p, spyIntoString(&out), spyIntoString(&errOut))
	if status != 0 {
		t.Error("unexpected status")
	}
	if errOut != "" {
		t.Error("should not have panicked")
	}
	r := regexp.MustCompile(expRegex)
	if !r.MatchString(out) {
		t.Errorf("expected pattern \"%s\" but got result \"%s\"", expRegex, out)
	}
}

func assertStatus(p parameters, exp int, t *testing.T) {
	status := run(p, func(s string, i ...interface{}) {}, func(s string, i ...interface{}) {})
	if status != exp {
		t.Errorf("expected status code %d but got %d", exp, status)
	}
}

func TestNowEncoding(t *testing.T) {
	assertSuccess(parameters{now: true, flagsSet: 1}, "^\\d{5}[[:alnum:]]{5,}$", t)
}

func TestDateEncoding(t *testing.T) {
	assertSuccess(parameters{date: "20190907134318", flagsSet: 1}, "^68495LTTOD$", t)
}

func TestNumberEncoding(t *testing.T) {
	assertSuccess(parameters{number: 4133980800, flagsSet: 1}, "^52247CRMTY$", t)
}

func TestBitstringEncoding(t *testing.T) {
	assertSuccess(parameters{bitstring: "01011101 01110011 10010111 11010110", flagsSet: 1}, "^68495LTTOD$", t)
}

func TestVerificationFull(t *testing.T) {
	assertStatus(parameters{reverse: "68495LTTOD", flagsSet: 1}, 0, t)
}

func TestVerificationPartial(t *testing.T) {
	assertStatus(parameters{reverse: "684", flagsSet: 1}, 4, t)
}

func TestVerificationBad(t *testing.T) {
	assertStatus(parameters{reverse: "B4D1NPUT", flagsSet: 1}, 1, t)
}

func TestBadUsageTooManyArguments(t *testing.T) {
	assertStatus(parameters{number: 42, reverse: "FOOBAR", flagsSet: 2}, 2, t)
}

func TestBadUsagePositionalArgument(t *testing.T) {
	assertStatus(parameters{number: 42, flagsSet: 1, leftoverArgs: true}, 2, t)
}

func TestBadUsageNoFlags(t *testing.T) {
	assertStatus(parameters{flagsSet: 0}, 2, t)
}

func TestBadUsageBitstringTooLong(t *testing.T) {
	assertStatus(parameters{bitstring: "11111111111111111111111111111111111111111111111111111111111111111", flagsSet: 1}, 2, t)
}

func TestBadUsageDateFormatInvalid(t *testing.T) {
	assertStatus(parameters{date: "2021", flagsSet: 1}, 2, t)
}
