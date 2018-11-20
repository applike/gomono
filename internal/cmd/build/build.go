package build

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/applike/gomono/internal/build"
	"github.com/applike/gomono/internal/build/golang"
	"github.com/applike/gomono/internal/build/makefile"
	"github.com/applike/gomono/internal/cmd"
	"github.com/applike/gomono/internal/dep"
	"github.com/applike/gomono/internal/dot"
	"github.com/applike/gomono/internal/search"
	"github.com/applike/gomono/internal/vcs"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/traverse"
)

var CmdBuild = &cmd.Command{
	Usage: "gomono build",
	Short: "Tool that builds what changed",
	Long: `
Long description of the gomono build tool
`,
}

var (
	print   = CmdBuild.Flag.Bool("print", false, "Print dependency graph and exit")
	all     = CmdBuild.Flag.Bool("all", false, "Build all, regardless of changes")
	from    = CmdBuild.Flag.String("from", "HEAD~1", "First commit")
	to      = CmdBuild.Flag.String("to", "HEAD", "Last commit")
	builder = CmdBuild.Flag.String("builder", "golang", "Build system to use. Possible values: 'golang', 'makefile'")
	action  = CmdBuild.Flag.String("action", "build", "Action to execute. Possible values: 'build', 'test'")
)

func init() {
	CmdBuild.Run = RunCmdBuild // break init cycle Gomono-RunCmd-flags-Gomono
}

func RunCmdBuild(cmd *cmd.Command, args []string) {
	var (
		mains []string
		err   error
	)

	switch *action {
	case "deploy", "build":
		mains, err = search.MainPackageNames(args)
	case "test":
		mains, err = search.PackageNames(args)
	default:
		log.Fatalf("unknown action %s", *action)
	}

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

		if *all {
			execute(nodes[m])
		} else {

			// Traverse graph, identify changes, build if something changed

			changedDeps, err := dep.DiffPkgs(nodes[m].Dir, *from, *to)
			if err != nil {
				log.Printf("failed to find changed Deps for %s, %v", m, err)
				execute(nodes[m])
				return
			}
			changedDepsMap := make(map[string]struct{})
			for _, d := range changedDeps {
				changedDepsMap[d] = struct{}{}
			}

			t := &traverse.BreadthFirst{
				EdgeFilter: func(_ graph.Edge) bool { return true },
			}
			n := t.Walk(g, nodes[m].Node, func(n graph.Node, _ int) bool {
				name := names[int(n.ID())]
				pkg := nodes[name]

				vendorName := ""
				if strings.Contains(pkg.Dir, "/vendor/") {
					vendorName = strings.Split(pkg.Dir, "/vendor/")[1]
				}
				if _, exists := changedDepsMap[vendorName]; exists || vcs.Changed(filepath.Join(pkg.Dir, "*.go"), *from, *to) {
					log.Printf("rebuilding %s because package in directory %v changed", m, pkg.Dir)
					return true
				}
				return false
			})
			if n != nil {
				execute(nodes[m])
			}
		}
	}
}

func execute(pkg *search.Package) {
	var (
		b   build.Builder
		err error
	)
	if *builder == "golang" {
		b, err = golang.NewFromImportPath(pkg.ImportPath)
		if err != nil {
			log.Fatalf("could not prepare build with arg %s: %v", pkg.ImportPath, err)
		}
	} else if *builder == "makefile" {
		b = makefile.NewFromDirectory(pkg.Dir)
	}

	switch *action {
	case "build":
		err = b.Build()
	case "test":
		err = b.Test()
	case "deploy":
		err = b.Deploy()
	default:
		log.Fatalf("unknown action %s", *action)
	}

	if err != nil {
		log.Fatalf("failed to build %s: %v", pkg.Dir, err)
	}
}
