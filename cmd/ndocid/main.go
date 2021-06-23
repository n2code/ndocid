package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/n2code/ndocid"
)

func main() {
	status := run(getParametersFromFlags(), sysOut, sysErrLineOut)
	os.Exit(status)
}

func sysOut(format string, msg ...interface{}) {
	fmt.Fprintf(os.Stdout, format, msg...)
}

func sysErrLineOut(format string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", msg...)
}

func getParametersFromFlags() (params parameters) {
	ref, _ := time.Parse(ndocid.DateFormat, ndocid.DateFormat)
	flag.StringVar(&params.date, "d", "", fmt.Sprintf("[MODE] Generate ID from given date and time.\nFor example `%s` which represents \"%s\".\nEvaluated in the machine's time zone.\nExit code greater than 0 if the input is not according to format.", ndocid.DateFormat, ref.Format("Mon Jan 2 15:04:05 2006")))
	flag.BoolVar(&params.now, "n", false, "[MODE] Generate ID from current date and time of this machine.")
	flag.StringVar(&params.bitstring, "b", "", "[MODE] Generate ID from string of bits, e.g. `\"00010110 11011011\"`.\nSpaces, tabs, underscores and leading zeros are being dropped.\nThe maximum length is 64 bits.\nBad input will result in an exit code greater than 0.")
	flag.Uint64Var(&params.number, "i", 0, "[MODE] Generate ID from number, e.g. `42`.\nAccepts any positive decimal number that can fit in an unsigned 64 bit integer.\nExit code greater than 0 if input exceeds range.")
	flag.BoolVar(&ndocid.Verbose, "v", false, "Verbose: More output.\nExplains algorithm in MODEs that generate IDs.\nProvides possible source representations when reversing is successful.")
	flag.StringVar(&params.reverse, "r", "", "[MODE] Reversing check: Validates given ID, e.g. `72639D77LD`.\nExit code 0: Valid full ID\nExit code 1: Invalid ID\nExit code 4: Plausible partial ID, needs further digits\nThe first line returned is OK / ERROR / PARTIAL for exit codes 0 / 1 / 4.")
	flag.Parse()
	params.flagsSet = flag.NFlag()
	if ndocid.Verbose {
		params.flagsSet-- //verbose does not count
	}
	if flag.NArg() != 0 {
		params.leftoverArgs = true
	}
	return
}
