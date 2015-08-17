package main

import (
	"fmt"
	"gongs/argparser"
	"gongs/biofile/fastq"
	"os"
	"runtime"
	"sync"
)

const countName = "count"
const countDesc = "count reads from fastq file"

var countArger = argparser.New(mainName, countName)

func countRunner(filenames ...string) {
	if len(filenames) == 0 {
		fmt.Fprintln(os.Stderr, mainName, countName, ": no input given!")
		os.Exit(1)
	}

	// setting multi-threads
	runtime.GOMAXPROCS(runtime.NumCPU())

	fqfiles := make([]*fastq.FastqFile, len(filenames))
	for i, filename := range filenames {
		fqfile, err := fastq.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fqfiles[i] = fqfile
	}

	chCount := make(chan int, len(filenames))
	go func(chCount chan int, fqfiles []*fastq.FastqFile) {
		wg := &sync.WaitGroup{}
		for _, fqfile := range fqfiles {
			wg.Add(1)
			go func(wg *sync.WaitGroup, fqfile *fastq.FastqFile, chCount chan int) {
				count := 0
				for fqfile.Next() {
					count++
				}
				chCount <- count
				fmt.Println(fqfile.Name, "contain Reads:", count)
				wg.Done()
			}(wg, fqfile, chCount)
		}
		wg.Wait()
		close(chCount)
	}(chCount, fqfiles)

	var tot int
	for c := range chCount {
		tot += c
	}
	fmt.Println("Total Reads:", tot)
}
