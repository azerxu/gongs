package main

import (
	"fmt"
	"gongs/argparser"
	"os"
)

const statName = "stat"
const statDesc = "stat fastq file"
const statVersion = "2015.08.19.1"

var statArger = argparser.New(mainName, statName)

func init() {
	statArger.Add("prefix", "-p", "--prefix", "output file prefix name", "stat")
	statArger.Add("thread", "-t", "--thread", "runtime thread Number default (all available cpu)", 0)
}

func statUsage() {
	statArger.Usage()
}

func statRunner(args ...string) {
	if len(args) == 0 {
		statUsage()
		os.Exit(1)
	}
	if err := statRun(args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func statRun(args ...string) error {
	if err := statArger.Parse(args...); err != nil {
		return err
	}

	// prefix := statArger.Get("prefix").(string)

	// setting multi-threads
	setThread(statArger.Get("thread").(int))
	return nil
}
