// Package command for command
package command

import "os"
import "fmt"
import "strings"

const sepSpace = 10

// SubCommand subcommand
type SubCommand struct {
	Name   string          // command name
	Desc   string          // command description
	Usage  func()          // show command usage
	Runner func(...string) // running command
}

// Command command
type Command struct {
	Name        string
	Desc        string
	Version     string
	subCommands map[string]*SubCommand
}

// Add add a subcmd to command
func (cmd *Command) Add(subCmd *SubCommand) error {
	_, ok := cmd.subCommands[subCmd.Name]
	if ok {
		return fmt.Errorf("SubCommand Name: %s has added", subCmd.Name)
	}
	cmd.subCommands[subCmd.Name] = subCmd
	return nil
}

// AddNew add a new subcommand
func (cmd *Command) AddNew(name string, desc string, usage func(), runner func(...string)) error {
	_, ok := cmd.subCommands[name]
	if ok {
		return fmt.Errorf("SubCommand Name: %s has added", name)
	}
	cmd.subCommands[name] = &SubCommand{Name: name, Desc: desc, Usage: usage, Runner: runner}
	return nil
}

// Usage  show command usage
func (cmd *Command) Usage() {
	l := sepSpace
	for cmdName := range cmd.subCommands {
		if len(cmdName) >= l {
			l = len(cmdName) + 1
		}
	}
	fmt.Fprintf(os.Stderr, "%s: %s (Version: %s)\n", cmd.Name, cmd.Desc, cmd.Version)
	fmt.Fprintf(os.Stderr, "Usage: %s <command> [options]\n", cmd.Name)
	fmt.Fprintln(os.Stderr, "  Command:")
	fmt.Fprintln(os.Stderr, "\t help", strings.Repeat(" ", l-len("help")), "show command help message")
	for subCmdName, subCmd := range cmd.subCommands {
		fmt.Fprintln(os.Stderr, "\t", subCmdName, strings.Repeat(" ", l-len(subCmdName)), subCmd.Desc)
	}
}

// HelpUsage show help usage
func (cmd *Command) HelpUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s help <command>\n", cmd.Name)
	fmt.Fprint(os.Stderr, "   Available Commands: help")
	for subCmdName := range cmd.subCommands {
		fmt.Fprint(os.Stderr, " ", subCmdName)
	}
	fmt.Fprintln(os.Stderr, "")
}

// Help show help
func (cmd *Command) Help(args ...string) {
	if len(args) == 0 {
		cmd.HelpUsage()
		return
	}

	subCmdName := args[0]
	if subCmdName == "help" {
		cmd.HelpUsage()
		return
	}

	subCmd, ok := cmd.subCommands[subCmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unkown SubCommand: %s\n\n", subCmdName)
		cmd.HelpUsage()
		return
	}
	subCmd.Usage()
}

// Run run command
func (cmd *Command) Run(args ...string) {
	if len(args) == 0 {
		cmd.Usage()
		return
	}

	subCmdName := args[0]
	args = args[1:]

	if subCmdName == "help" {
		cmd.Help(args...)
		return
	}
	subCmd, ok := cmd.subCommands[subCmdName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unkown SubCommand: %s\n\n", subCmdName)
		cmd.Usage()
		return
	}
	subCmd.Runner(args...)
}

// New New(name, [desc, version])
func New(args ...string) *Command {
	name := "command"
	desc := "a command tool kit"
	version := "1.0"

	l := len(args)
	switch {
	case l >= 3:
		version = args[2]
		fallthrough
	case l == 2:
		desc = args[1]
		fallthrough
	case l == 1:
		name = args[0]
	}
	return &Command{
		Name:        name,
		Desc:        desc,
		Version:     version,
		subCommands: make(map[string]*SubCommand),
	}
}
