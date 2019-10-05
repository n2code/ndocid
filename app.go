package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	benchmark()

	status := run(getParametersFromFlags(), sysOut, sysErrOut)

	os.Exit(status)
}

func sysOut(format string, msg ...interface{}) {
	fmt.Fprintf(os.Stdout, format, msg...)
}

func sysErrOut(format string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, format, msg...)
}

func getParametersFromFlags() (params parameters) {
	ref, _ := time.Parse(dateFormat, dateFormat)
	flag.StringVar(&params.date, "d", "", fmt.Sprintf("[MODE] Generate ID from given date and time.\nFor example `%s` which represents \"%s\".\nEvaluated in the machine's time zone.", dateFormat, ref.Format("Mon Jan 2 15:04:05 2006")))
	flag.BoolVar(&params.now, "n", false, "[MODE] Generate ID from current date and time of this machine.")
	flag.StringVar(&params.bitstring, "b", "", "[MODE] Generate ID from string of bits, e.g. `\"00010110 11011011\"`.\nSpaces, tabs, underscores and leading zeros are being dropped.\nThe maximum length is 64 bits.")
	flag.Uint64Var(&params.number, "i", 0, "[MODE] Generate ID from number, e.g. `42`.\nAccepts any positive decimal number that can fit in an unsigned 64 bit integer.")
	flag.BoolVar(&verbose, "v", false, "Verbose: More output.\nExplains algorithm in MODEs that generate IDs.\nProvides possible source representations when reversing is successful.")
	flag.StringVar(&params.reverse, "r", "", "[MODE] Reversing check: Validates given ID, e.g. `72639D77LD`.\nExit code 0: Valid full ID\nExit code 1: Invalid ID\nExit code 4: Plausible partial ID, needs further digits\nThe first line returned is OK / ERROR / PARTIAL for exit codes 0 / 1 / 4.")
	flag.Parse()
	params.flagsSet = flag.NFlag()
	if verbose {
		params.flagsSet-- //verbose does not count
	}
	if flag.NArg() != 0 {
		params.leftoverArgs = true
	}
	return
}

func benchmark() {
	switch len(os.Args) {
	case 2:
		switch os.Args[1] {
		case "iterate":
			for i := uint64(0); i < math.MaxUint64; i++ {
				fmt.Println(encode(i))
			}
		}
	case 5:
		switch os.Args[1] {
		case "distribution":
			from, errFrom := strconv.ParseUint(os.Args[3], 10, 64)
			if errFrom != nil {
				sysErrOut("%s", errFrom)
				os.Exit(1)
			}
			to, errTo := strconv.ParseUint(os.Args[4], 10, 64)
			if errTo != nil {
				sysErrOut("%s", errTo)
				os.Exit(1)
			}

			prefixes1 := make(map[string]int, 8)
			prefixes2 := make(map[string]int, 8*8)
			prefixes3 := make(map[string]int, 8*8*8)

			recordedEncodeOf := func(in uint64) {
				encoded := encode(in)
				prefixes1[encoded[:1]]++
				prefixes2[encoded[:2]]++
				prefixes3[encoded[:3]]++
			}

			switch os.Args[2] {
			case "all":
				fmt.Println("Calculating full prefix distribution...")
				feedback := time.Tick(5000 * time.Millisecond)
				for i := from; i <= to; i++ {
					select {
					case <-feedback:
						fmt.Printf("Currently at %d...\n", i)
					default:
						recordedEncodeOf(i)
					}
				}
			case "random5percent":
				fmt.Println("Calculating random prefix distribution...")
				feedback := time.Tick(5000 * time.Millisecond)
				for i := from; i <= from+(to-from+1)/20; i++ {
					select {
					case <-feedback:
						fmt.Println("Still generating random IDs in range...")
					default:
						recordedEncodeOf(from + uint64(rand.Int63n(int64(to-from+1))))
					}
				}
			}

			printPrefixes := func(in map[string]int) {
				raw := strings.ReplaceAll(fmt.Sprint(in), " ", "\n")
				fmt.Println(raw[4 : len(raw)-1])
			}
			printPrefixes(prefixes1)
			printPrefixes(prefixes2)
			printPrefixes(prefixes3)
			os.Exit(0)
		}
	}
}
