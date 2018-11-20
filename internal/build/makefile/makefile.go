package makefile

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Makefile struct {
	file string
}

func NewFromDirectory(dir string) *Makefile {
	return &Makefile{
		file: filepath.Join(dir, "Makefile"),
	}
}

func (m *Makefile) Build() error {
	return runMake(m.file, []string{"build"})
}
func (m *Makefile) Test() error {
	return runMake(m.file, []string{"test"})
}

func (m *Makefile) Deploy() error {
	return runMake(m.file, []string{"deploy"})
}

func runMake(f string, args []string) error {
	cmd := exec.Command("make", args...)
	cmd.Dir = filepath.Dir(f)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
