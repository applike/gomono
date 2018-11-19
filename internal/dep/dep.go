package dep

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
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

// Projects contains a list of projects listed in deps Gopkg.lock
type Projects struct {
	Projects []project
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

// Diff returns a list of projets, which are different in a and b
func Diff(old, new *Projects) *Projects {
	var changed Projects
	for _, p := range new.Projects {
		newP, found := find(p, old.Projects)
		if !found || newP.Version != p.Version || newP.Revision != p.Revision {
			changed.Projects = append(changed.Projects, p)
		}
	}
	return &changed
}
