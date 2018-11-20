package golang

import (
	"os"
	"os/exec"
)

type Golang struct {
	dir string
	out string
}

func (g *Golang) Build() error {
	cmd := exec.Command("go", "build", "-o", g.out, g.dir)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (g *Golang) Test() error {
	cmd := exec.Command("go", "test", g.dir)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
