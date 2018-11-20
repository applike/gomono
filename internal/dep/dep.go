package dep

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/applike/gomono/internal/vcs"
)

func getLockFile(p string) (string, error) {
	path, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	if path == "/" {
		return "", fmt.Errorf("could not find project Gopkg.lock")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if f.Name() == "Gopkg.lock" {
			return filepath.Join(path, f.Name()), nil
		}
	}

	return getLockFile(filepath.Dir(path))
}

type pkg string

type project struct {
	Digest    string
	Name      string
	Pruneopts string
	Revision  string
	Version   string

	Packages []pkg
}

func pkgEqual(a, b []pkg) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (p *project) equals(o *project) bool {
	return p.Name == o.Name &&
		p.Digest == o.Digest &&
		p.Revision == o.Revision &&
		p.Version == o.Version &&
		pkgEqual(p.Packages, o.Packages)
}

// Projects contains a list of projects listed in deps Gopkg.lock
type Projects struct {
	Projects []*project
}

func parseTOML(input string) (*Projects, error) {
	var p Projects
	_, err := toml.Decode(input, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func find(p project, l []project) (project, bool) {
	for _, lp := range l {
		if lp.Name == p.Name {
			return lp, true
		}
	}
	return p, false
}

// DiffPkgs returns a list of packages which have been
// changed between the commits named by old and new
func DiffPkgs(path, old, new string) ([]string, error) {
	lock, err := getLockFile(path)
	if err != nil {
		return nil, err
	}

	o, err := vcs.Show(lock, old)
	if err != nil {
		return nil, err
	}

	n, err := vcs.Show(lock, new)
	if err != nil {
		return nil, err
	}

	op, err := parseTOML(o)
	if err != nil {
		return nil, err
	}
	np, err := parseTOML(n)
	if err != nil {
		return nil, err
	}

	var pkgs = make([]string, 0)
	changed := diff(op.Projects, np.Projects)
	for _, p := range changed {
		if len(p.Packages) <= 0 {
			pkgs = append(pkgs, string(p.Name))
		}
		for _, pkg := range p.Packages {
			pkgs = append(pkgs, fmt.Sprintf("%s/%s", p.Name, string(pkg)))
		}
	}

	return pkgs, nil
}

func diff(a, b []*project) []*project {
	var diff []*project

	for i := 0; i < 2; i++ {
		for _, s1 := range a {
			found := false
			for _, s2 := range b {
				if s1.equals(s2) {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			a, b = b, a
		}
	}
	return diff
}
