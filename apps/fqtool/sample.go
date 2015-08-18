package main

import (
	"bufio"
	"fmt"
	"gongs/argparser"
	"gongs/biofile/fastq"
	"gongs/xopen"
	"math/rand"
	"os"
	"runtime"
	"time"
)

const sampleName = "sample"
const sampleDesc = "sample a sub set from fastq file"

var sampleArger = argparser.New(mainName, sampleName)

func init() {
	sampleArger.Add("rate", "-r", "--rate", "sample rate", 0.1)
	sampleArger.Add("single", "-s", "--single", "input file is single file", false)
	sampleArger.Add("prefix", "-p", "--prefix", "output file prefix name", "sample")
	sampleArger.Add("seed", "-S", "--seed", "random seed", 0)
	sampleArger.Add("thread", "-t", "--threads", "threads number default use all", 0)
}

func sampleRunner(args ...string) {
	if len(args) == 0 {
		sampleArger.Usage()
		os.Exit(1)
	}

	if err := sampleArger.Parse(args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rate := sampleArger.Get("rate").(float64)
	single := sampleArger.Get("single").(bool)
	prefix := sampleArger.Get("prefix").(string)
	seed := int64(sampleArger.Get("seed").(int))
	thread := sampleArger.Get("thread").(int)

	if cpus := runtime.NumCPU(); thread < 1 || thread > cpus {
		thread = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(thread)

	if seed != 0 {
		rand.Seed(seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	if err := sampleRun(single, rate, prefix, sampleArger.Args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func sampleRun(single bool, rate float64, prefix string, filenames ...string) error {
	if single {
		return sampleSingleRun(rate, prefix, filenames...)
	}
	return samplePairRun(rate, prefix, filenames...)
}

func sampleSingleRun(rate float64, prefix string, filenames ...string) error {
	out, err := xopen.Xcreate(prefix+".fastq", "w")
	if err != nil {
		return err
	}
	defer out.Close()
	outter := bufio.NewWriter(out)
	defer outter.Flush()

	fqChan, errChan := fastq.Load(filenames...)
	for {
		select {
		case fq, ok := <-fqChan:
			if !ok {
				return nil
			}
			if rate > rand.Float64() {
				fmt.Fprintln(outter, fq)
			}
		case err := <-errChan:
			return err
		}
	}
	return nil
}

func samplePairRun(rate float64, prefix string, filenames ...string) error {
	out1, err := xopen.Xcreate(prefix+".r1.fastq", "w")
	if err != nil {
		return err
	}
	defer out1.Close()
	out2, err := xopen.Xcreate(prefix+".r2.fastq", "w")
	if err != nil {
		return err
	}
	defer out2.Close()
	outter1 := bufio.NewWriter(out1)
	outter2 := bufio.NewWriter(out2)
	defer outter1.Flush()
	defer outter2.Flush()

	pchan, errchan := fastq.LoadPair(filenames...)
	for {
		select {
		case p, ok := <-pchan:
			if !ok { // pchan closed, all ok
				return nil
			}
			if rate > rand.Float64() {
				fmt.Fprintln(outter1, p.Read1)
				fmt.Fprintln(outter2, p.Read1)
			}
		case err := <-errchan: // something wrong at input fastq files
			return err
		}
	}
	return nil
}
