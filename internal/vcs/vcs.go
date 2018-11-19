package vcs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Changed returns whether a directory changed between the first and the second commit
func Changed(dir, first, second string) bool {
	var (
		cmdOut []byte
		err    error
	)

	cmd := exec.Command("git", "diff", "--name-only", first, second, dir)
	if cmdOut, err = cmd.CombinedOutput(); err != nil {
		return true
	}

	return len(strings.TrimSpace(string(cmdOut))) > 0
}

// Show returns the lines of a file at revision rev
func Show(path, rev string) (string, error) {

	var (
		cmdOut []byte
		err    error
	)

	rel, err := absFromRoot(path)
	if err != nil {
		return "", err
	}
	root, err := root()
	if err != nil {
		return "", err
	}
	git := filepath.Join(root, ".git")
	cmd := exec.Command(
		"git",
		fmt.Sprintf("%s%s", "--work-tree=", root),
		fmt.Sprintf("%s%s", "--git-dir=", git),
		"show",
		fmt.Sprintf("%s:%s", rev, rel),
	)
	cmd.Stderr = os.Stderr
	if cmdOut, err = cmd.Output(); err != nil {
		return "", err
	}

	return string(cmdOut), nil
}

// AbsFromRoot returns the path of a file in a git repository
// with the root path of that git directory removed
func absFromRoot(p string) (string, error) {
	root, err := root()
	if err != nil {
		return "", err
	}

	return p[len(root)+1:], nil
}

func root() (string, error) {
	var (
		cmdOut []byte
		err    error
	)
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stderr = os.Stderr
	if cmdOut, err = cmd.Output(); err != nil {
		return "", err
	}

	return strings.TrimSpace(string(cmdOut)), nil
}
