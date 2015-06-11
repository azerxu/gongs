package main

import (
	"fmt"
	"gongs/arger"
	"gongs/biofile/fastq"
	"gongs/lib"
	"math/rand"
	"os"
	"time"
)

const sampleName = "sample"
const sampleDesc = "sample a sub set from fastq file"

var sampleArger = arger.New(mainName, sampleName)

func init() {
	sampleArger.Add("rate", "-r", "--rate", "sample rate", 0.1)
	sampleArger.Add("single", "-s", "--single", "input file is single file", false)
	sampleArger.Add("prefix", "-p", "--prefix", "output file prefix name", "sample")
	sampleArger.Add("seed", "-S", "--seed", "random seed", 0)
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
	fqchan, err := fastq.Loads(filenames...)
	if err != nil {
		return err
	}

	outter, err := lib.Xcreate(prefix+".fastq", "w")
	if err != nil {
		return err
	}
	defer outter.Close()

	for fq := range fqchan {
		if rate > rand.Float64() {
			fmt.Fprintln(outter, fq)
		}
	}
	return nil
}

func samplePairRun(rate float64, prefix string, filenames ...string) error {
	pairchan, err := fastq.LoadPairs(filenames...)
	if err != nil {
		return err
	}
	outter1, err := lib.Xcreate(prefix+".r1.fastq", "w")
	if err != nil {
		return err
	}
	defer outter1.Close()
	outter2, err := lib.Xcreate(prefix+".r2.fastq", "w")
	if err != nil {
		return err
	}
	defer outter2.Close()
	for p := range pairchan {
		if rate > rand.Float64() {
			fmt.Fprintln(outter1, p.Read1)
			fmt.Fprintln(outter2, p.Read2)
		}
	}
	return nil
}
