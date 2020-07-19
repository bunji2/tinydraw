package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	usage = "Usage: %s file.js [ grid_width grid_height margin ]"
)

func main() {
	os.Exit(run())
}

func run() int {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		return 1
	}
	jsFilePath := os.Args[1]

	if len(os.Args) > 4 {
		tmp, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		gw := tmp
		tmp, err = strconv.Atoi(os.Args[3])
		if err != nil {
			panic(err)
		}
		gh := tmp
		m := 0
		if len(os.Args) > 4 {
			tmp, err = strconv.Atoi(os.Args[4])
			if err != nil {
				panic(err)
			}
			m = tmp
		}
		SetParams(gw, gh, m)
	}

	err := RunJS(jsFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	return 0
}
