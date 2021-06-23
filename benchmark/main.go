package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/n2code/ndocid"
)

func main() {
	switch len(os.Args) {
	case 2:
		switch os.Args[1] {
		case "iterate":
			for i := uint64(0); i < math.MaxUint64; i++ {
				fmt.Println(ndocid.EncodeUint64(i))
			}
		default:
			exitWithUsageError("Bad argument!")
		}
		os.Exit(0)
	case 5:
		switch os.Args[1] {
		case "distribution":
			from, errFrom := strconv.ParseUint(os.Args[3], 10, 64)
			if errFrom != nil {
				sysErrLineOut(errFrom)
				os.Exit(1)
			}
			to, errTo := strconv.ParseUint(os.Args[4], 10, 64)
			if errTo != nil {
				sysErrLineOut(errTo)
				os.Exit(1)
			}

			prefixes1 := make(map[string]int, 8)
			prefixes2 := make(map[string]int, 8*8)
			prefixes3 := make(map[string]int, 8*8*8)

			recordedEncodeOf := func(in uint64) {
				encoded := ndocid.EncodeUint64(in)
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
			default:
				exitWithUsageError("Bad argument!")
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
	default:
		exitWithUsageError("Bad number of arguments!")
	}
}

func sysErrLineOut(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
}

func exitWithUsageError(msg interface{}) {
	sysErrLineOut(msg)
	sysErrLineOut("\nUsage: benchmark iterate\n    OR benchmark distribution all|random5percent FROM TO\n       (e.g. benchmark distribution all 1000000 9999999)")
	os.Exit(2)
}
