package search

import (
	"go/build"
	"os"
	"os/exec"
	"strings"
)

func Import(path string) (*build.Package, error) {
	var (
		p   *build.Package
		err error
		cwd string
	)
	cwd, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	p, err = build.Import(path, cwd, 0)
	return p, err
}

// Packages returns a list of all packages named by pattern
func Packages(pattern string) ([]*build.Package, error) {
	paths, err := ImportPaths(pattern)
	if err != nil {
		return nil, err
	}

	var pkgs = make([]*build.Package, 0)
	for _, p := range paths {
		pkg, err := Import(p)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}

// ImportPaths returns the ImportPaths of all packages named by pattern
func ImportPaths(pattern string) ([]string, error) {
	return List([]string{"-e", pattern})
}

func List(args []string) ([]string, error) {
	subcmd := []string{
		"list",
	}
	args = append(subcmd, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	cmdOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	pkgs := strings.Split(string(cmdOut), "\n")
	var packages = []string{}
	for _, p := range pkgs {
		if len(p) > 0 {
			packages = append(packages, p)
		}
	}
	return packages, nil
}

func MainPackages(pattern string) ([]string, error) {

	args := []string{
		"-f",
		"{{if eq .Name \"main\" }}{{ .ImportPath }}{{end}}",
		pattern,
	}
	return List(args)
}
