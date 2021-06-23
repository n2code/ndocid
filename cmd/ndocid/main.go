package main

import (
	"flag"
	"fmt"
	"os"

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
	flag.StringVar(&params.date, "d", "", "DATE-MODE: Generate ID from given date and time.\n  For example `20060102150405` which represents \"Mon Jan 2 15:04:05 2006\".\n  Evaluated in the machine's time zone.\n  Exit code greater than 0 if the input is not according to format.")
	flag.BoolVar(&params.now, "n", false, "NOW-MODE: Generate ID from current date and time of this machine.")
	flag.StringVar(&params.bitstring, "b", "", "BITSTRING-MODE: Generate ID from string of bits, e.g. `\"00010110 11011011\"`.\n  Spaces, tabs, underscores and leading zeros are being dropped.\n  The maximum length is 64 bits.\n  Bad input will result in an exit code greater than 0.")
	flag.Uint64Var(&params.number, "i", 0, "INTEGER-MODE: Generate ID from number, e.g. `42`.\n  Accepts any positive decimal number that can fit in an unsigned 64 bit integer.\n  Exit code greater than 0 if input exceeds range.")
	flag.BoolVar(&ndocid.Verbose, "v", false, "Verbose option: Generate more human-readable output.\n  Explains algorithm in MODEs that generate IDs.\n  Provides possible source representations when reversing is successful.")
	flag.StringVar(&params.reverse, "r", "", "REVERSING/CHECK-MODE: Validates given ID, e.g. `72639D77LD`.\n  Exit code 0: Valid full ID\n  Exit code 1: Invalid ID\n  Exit code 4: Plausible partial ID (beginning), needs further digits\n  The first line returned is OK / ERROR / PARTIAL for exit codes 0 / 1 / 4.")
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
