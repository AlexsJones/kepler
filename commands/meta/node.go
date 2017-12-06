// Package meta is aimed at working
// within a git repo that manages other repos.
// The issues this package helps solve is that
// node packages are a pain with npm as you
// can't say use local if you can find them
// instead of using remote.
package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	git "gopkg.in/src-d/go-git.v4"

	Node "github.com/AlexsJones/kepler/commands/node"
	"github.com/AlexsJones/kepler/commands/submodules"
)

// Information ...
type Information struct {
	Projects map[string]*Node.PackageJSON
}

// NewInformation creates a struct containing information about the meta repo
func NewInformation() (*Information, error) {
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		return nil, errors.New("Can not create information about meta repo")
	}
	data := &Information{
		Projects: make(map[string]*Node.PackageJSON),
	}
	submodules.LoopSubmodules(func(sub *git.Submodule) {
		filepath := path.Join(sub.Config().Path, "package.json")
		if _, node := os.Stat(filepath); !os.IsNotExist(node) {
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				return
			}
			var p Node.PackageJSON
			json.Unmarshal(b, &p)
			data.Projects[sub.Config().Name] = &p
		}
	})
	return data, nil
}

// ResolveLocalDependancies will explore (via some graph expansion)
// once it is completed, it will return the list of the required
// pacakages otherwise, return an informative error
func (meta *Information) ResolveLocalDependancies(project string) ([]string, error) {
	if _, exists := meta.Projects[project]; !exists {
		return nil, fmt.Errorf("%s does not exists", project)
	}
	ResolvedDeps := make(map[string]bool)
	Explore := make(map[string]*Node.PackageJSON)
	// Making sure we don't try to explore the started node
	// if it is required by another project
	ResolvedDeps[project] = true
	Explore[project] = meta.Projects[project]
	for len(Explore) > 0 {
		for node, pack := range Explore {
			for name := range pack.Dependencies {
				if ResolvedDeps[name] {
					// Nothing to do as its already been resolved
					continue
				}
				if _, local := meta.Projects[name]; local {
					Explore[name] = meta.Projects[name]
				}
			}
			for name := range pack.DevDependencies {
				if ResolvedDeps[name] {
					// Nothing to do as its already been resolved
					continue
				}
				if _, local := meta.Projects[name]; local {
					Explore[name] = meta.Projects[name]
				}
			}
			ResolvedDeps[node] = true
			delete(Explore, node)
		}
	}
	// Make sure we don't include ourselves when we print out
	delete(ResolvedDeps, project)
	deps := []string{}
	for dep := range ResolvedDeps {
		deps = append(deps, dep)
	}
	return deps, nil
}
