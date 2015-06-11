package main

import (
	"fmt"
	"gongs/arger"
	"gongs/biofile/fastq"
	"os"
)

const countName = "count"
const countDesc = "count reads from fastq file"

var countArger = arger.New(mainName, countName)

func countRunner(args ...string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, mainName, countName, ": no input given!")
		os.Exit(1)
	}
	for _, file := range args {
		fmt.Println("file:", file)
		fqfile, err := fastq.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		cc := 0
		for _ = range fqfile.Load() {
			cc++
		}
		fmt.Println(file, "contain Reads:", cc)
	}
}
