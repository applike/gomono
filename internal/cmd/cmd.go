package cmd

import (
	"flag"
	"strings"
)

type Command struct {
	Run      func(cmd *Command, args []string)
	Flag     flag.FlagSet
	Usage    string
	Short    string
	Long     string
	Commands []*Command
}

var Gomono = &Command{
	Usage: "gomono",
	Short: "Tool that tells what changed",
	Long: `
Long description of the gomono tool
`,
}

func (c *Command) LongName() string {
	name := c.Usage
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	}

	return name
}

func (c *Command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}
