package cmd

import (
	"flag"
	"strings"
)

// Command stores the information about a runnable action
type Command struct {
	Run      func(cmd *Command, args []string)
	Flag     flag.FlagSet
	Usage    string
	Short    string
	Long     string
	Commands []*Command
}

// Gomono is the main command
var Gomono = &Command{
	Usage: "gomono [cmd] [flags] [packages]",
	Short: "gomono is a tool for analyzing changes in source code",
	Long: `

`,
}

// LongName returns the long name of a command
func (c *Command) LongName() string {
	name := c.Usage
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	}

	return name
}

// Name returns the name of a command
func (c *Command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}
