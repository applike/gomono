package build

import (
	"log"
	"os"

	"github.com/applike/gomono/internal/cmd"
	"github.com/applike/gomono/internal/dot"
	"github.com/applike/gomono/internal/search"
)

var CmdBuild = &cmd.Command{
	Usage: "gomono build",
	Short: "Tool that builds what changed",
	Long: `
Long description of the gomono build tool
`,
}

var (
	print = CmdBuild.Flag.Bool("print", false, "Print dependency graph and exit")
)

func init() {
	CmdBuild.Run = RunCmdBuild // break init cycle Gomono-RunCmd-flags-Gomono
}

func RunCmdBuild(cmd *cmd.Command, args []string) {

	mains, err := search.MainPackages(args)
	if err != nil {
		log.Fatalf("failed to get packages %s, %v", args, err)
	}

	for _, m := range mains {
		g, nodes, err := search.Graph([]string{m})
		if err != nil {
			log.Fatalf("failed to build graph for %s: %v", m, err)
		}
		names := make(map[int]string)
		for _, v := range nodes {
			names[int(v.ID())] = v.ImportPath
		}

		if *print {
			err = dot.Dot(os.Stdout, m, g, names)
			if err != nil {
				log.Fatalf("failed to print graph to stdout: %v", err)
			}
		}
	}
}
