package main

import (
	"fmt"
	"gongs/biofile/fastq"
	"os"
)

func main() {
	for _, file := range os.Args[1:] {
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
