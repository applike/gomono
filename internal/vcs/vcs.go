package vcs

import (
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

func depChanged(pkg, first, second string) bool {
	return false
}
