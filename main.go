package main

import (
	"flag"
	"log"

	"github.com/applike/gomono/internal/cmd"
	"github.com/applike/gomono/internal/cmd/build"
)

func main1() {

}

func init() {
	cmd.Gomono.Commands = []*cmd.Command{
		build.CmdBuild,
	}
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range cmd.Gomono.Commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			cmd.Run(cmd, args)
		}
	}
}

func usage() {
	log.Fatal(cmd.Gomono.Usage)
}
