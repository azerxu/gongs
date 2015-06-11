package main

import (
	"gongs/command"
	"os"
)

const mainName = "fqtool"
const mainDesc = "a tool kit for fastq format files"
const mainVersion = "2015.06.11.1"

func main() {
	cmd := command.New(mainName, mainDesc, mainVersion)
	cmd.Add(&command.SubCommand{ // add count command
		Name:   countName,
		Desc:   countDesc,
		Usage:  countArger.Usage,
		Runner: countRunner})
	cmd.Add(&command.SubCommand{ // add sample command
		Name:   sampleName,
		Desc:   countDesc,
		Usage:  sampleArger.Usage,
		Runner: sampleRunner})
	cmd.Run(os.Args[1:]...)
}
