package golang

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
)

type Golang struct {
	importPath string
	dir        string
	out        string
}

func NewFromImportPath(path string) (*Golang, error) {
	p, err := build.Import(path, "", 0)
	if err != nil {
		return nil, err
	}

	return &Golang{
		importPath: p.ImportPath,
		dir:        p.Dir,
		out:        p.PkgObj,
	}, nil
}

func (g *Golang) Build() error {
	out := g.out
	args := []string{"build"}
	if len(out) > 0 {
		args = append(args, "-o", g.out)
	}
	cmd := exec.Command("go", append(args, g.importPath)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (g *Golang) Test() error {
	cmd := exec.Command("go", "test", g.importPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (g *Golang) Deploy() error {
	return fmt.Errorf("not implemented")
}
