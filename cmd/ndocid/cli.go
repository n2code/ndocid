package main

import (
	"time"

	"github.com/n2code/ndocid"
)

type outFunc func(string, ...interface{})

var verboseLineOut = func(format string, msg ...interface{}) {
	if ndocid.Verbose {
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
	const seeUsage = "(see -h for usage)"
	switch {
	case p.flagsSet == 0:
		errOut(`No [MODE] flag given! %s`, seeUsage)
		return 1
	case p.flagsSet > 1:
		errOut(`Only one [MODE] flag may be set at a time %s`, seeUsage)
		return 1
	case p.leftoverArgs:
		errOut(`Leftover arguments after flags %s`, seeUsage)
		return 1
	}

	if p.reverse != "" {
		decoded, err, complete := ndocid.Decode(p.reverse)
		if err != nil {
			out("INVALID\n")
			errOut("%s", err)
			return 1
		}
		if complete {
			out("OK\n")
			verboseLineOut("Integer: %d", decoded)
			verboseLineOut("Date: %s", time.Unix(int64(decoded), 0).Format(time.RFC1123Z))
			verboseLineOut("Bitstring: %b", decoded)
		} else {
			out("PARTIAL\n")
			return 4
		}
	} else {
		var encoded string
		switch {
		case p.date != "":
			var err error
			encoded, err = ndocid.EncodeDatetime(p.date)
			if err != nil {
				errOut("%s", err)
				return 2
			}
		case p.now:
			rightNow := time.Now()
			verboseLineOut("Using current point in time: %s (unix time in seconds: %d)", rightNow.Format(time.RFC1123Z), rightNow.Unix())
			encoded = ndocid.EncodeUint64(uint64(rightNow.Unix()))
		case p.bitstring != "":
			var err error
			encoded, err = ndocid.EncodeBitstring(p.bitstring)
			if err != nil {
				errOut("%s", err)
				return 2
			}
		case p.reverse != "":
		default:
			encoded = ndocid.EncodeUint64(p.number)
		}

		verboseLineOut("Resulting encoded ID:")
		out(encoded)
	}
	return 0
}
