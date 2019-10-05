package main

import (
	"time"
)

type outFunc func(string, ...interface{})

var verbose bool
var verboseOut = func(format string, msg ...interface{}) {
	if verbose {
		sysOut(format+"\n", msg...)
	}
}

type parameters struct {
	bitstring    string
	date         string
	now          bool
	number       uint64
	reverse      string
	flagsSet     int
	leftoverArgs bool
}

func run(p parameters, out outFunc, errOut outFunc) (status int) {
	const seeUsage = " (see -h for usage)"
	switch {
	case p.flagsSet == 0:
		errOut(`No [MODE] flag given!`, seeUsage)
		return 1
	case p.flagsSet > 1:
		errOut(`Only one [MODE] flag may be set at a time`, seeUsage)
		return 1
	case p.leftoverArgs:
		errOut(`Leftover arguments after flags`, seeUsage)
		return 1
	}

	if p.reverse != "" {
		decoded, err, complete := decode(p.reverse)
		if err != nil {
			out("ERROR\n")
			errOut("%s", err)
			return 1
		}
		if complete {
			out("OK\n")
			verboseOut("Integer: %d", decoded)
			verboseOut("Date: %s", time.Unix(int64(decoded), 0).Format(time.RFC1123Z))
			verboseOut("Bitstring: %b", decoded)

		} else {
			out("PARTIAL\n")
		}
	} else {
		var encoded string
		switch {
		case p.date != "":
			var err error
			encoded, err = encodeDatetime(p.date)
			if err != nil {
				errOut("%s", err)
				return 1
			}
		case p.now:
			rightNow := time.Now()
			verboseOut("Using current point in time: %s (unix time in seconds: %d)", rightNow.Format(time.RFC1123Z), rightNow.Unix())
			encoded = encodeUint64(uint64(rightNow.Unix()))
		case p.bitstring != "":
			var err error
			encoded, err = encodeBitstring(p.bitstring)
			if err != nil {
				errOut("%s", err)
				return 1
			}
		case p.reverse != "":
		default:
			encoded = encodeUint64(p.number)
		}

		verboseOut("Resulting encoded ID:")
		out(encoded)
	}
	return 0
}
