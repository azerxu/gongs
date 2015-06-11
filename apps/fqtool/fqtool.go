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
	cmd.Add(&command.SubCommand{Name: countName, Desc: countDesc, Usage: countArger.Usage, Runner: countRunner})
	cmd.Run(os.Args[1:]...)
}
