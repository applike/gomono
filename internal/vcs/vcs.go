package vcs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Changed returns whether a package changed between the first and the second commit
func Changed(pkg, first, second string) bool {
	if strings.Contains(pkg, "/vendor/") {
		return false
	}
	return changed(pkg, first, second)
}

func changed(dir, first, second string) bool {
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
	cmd := exec.Command("git", "show", fmt.Sprintf("%v:%v", rev, path))
	cmd.Stderr = os.Stderr
	if cmdOut, err = cmd.Output(); err != nil {
		return "", err
	}

	return string(cmdOut), nil
}
