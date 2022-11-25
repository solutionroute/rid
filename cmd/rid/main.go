// Package main - the `rid` command - generate or inspect rid.
package main

import (
	"flag"
	"fmt"
	"os"
    "strings"

	"github.com/solutionroute/rid"
)

var (
	count   int  = 1
)

func init() {
	flag.IntVar(&count, "c", count, "Generate n * IDs")
}

func main() {
    flag.Usage = func() {
        pgm := os.Args[0]
        fmt.Fprintf(flag.CommandLine.Output(), "usage: %s -c N          # print N rid(s)\n", pgm)
        fmt.Fprintf(flag.CommandLine.Output(), "       %s 0629p0rqdrw8p # decode one or more rid(s)\n", pgm)
        // flag.PrintDefaults()
    }
	flag.Parse()
	args := flag.Args()
    
    if count > 1 && len(args) > 0 {
        fmt.Fprintf(flag.CommandLine.Output(), "error: -c (output) and args (input) both specified; perform only one at a time.\n")
        flag.Usage()
        os.Exit(1)
    }

	errors := 0

    // If no valid flag, Either attempt to decode string as a rid
	for _, arg := range args {
		id, err := rid.FromString(arg)
		if err != nil {
			errors++
			fmt.Printf("[%s] %s\n", arg, err)
			continue
		}
        fmt.Printf("[%s] seconds:%d random:%d machine:%v pid:%v time:%v ID{%s}\n", 
            arg, id.Seconds(), id.Random(), id.Machine(), id.Pid(), id.Time(), asHex(id[:]))
	}
	if errors > 0 {
		fmt.Printf("%d error(s)\n", errors)
		os.Exit(1)
	}

	// OR... generate one (or -c value) rid
    if len(args) == 0 {
        for c := 0; c < count; c++ {
            fmt.Fprintf(os.Stdout, "%s\n", rid.New())
        }

    }
}

func asHex(b []byte) string {
    s := []string{}
    for _, v := range b {
        s = append(s, fmt.Sprintf("%#x", v))
    }
    return strings.Join(s, ", ")

}
