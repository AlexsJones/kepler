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

	"github.com/Alexsjones/kepler/commands/types"
	"github.com/fatih/color"
)

// Information ...
type Information struct {
	// projects, poor mans version of a set
	Projects map[string]*types.PackageJSON
}

func loopSubmodules(path string, callback func(sub *git.Submodule) error) error {

	r, err := git.PlainOpen(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	worktree, err := r.Worktree()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	submodules, err := worktree.Submodules()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, sub := range submodules {
		callback(sub)
	}
	if len(submodules) == 0 {
		color.Red("No submodules found")
	}
	return nil
}

//LoopSubmodules will run through all submodules in the current repository
//It will return a nil error object on success
func LoopSubmodules(callback func(sub *git.Submodule)) error {
	loopSubmodules(".", func(sub *git.Submodule) error {
		return callback(sub)
	})
	return nil
}

// NewInformation creates a struct containing information about the meta repo
func NewInformation() (*Information, error) {
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		return nil, errors.New("Can not create information about meta repo")
	}
	data := &Information{}
	LoopSubmodules(func(sub *git.Submodule) {
		filepath := path.Join(sub.Config().Path, "package.json")
		if _, node := os.Stat(filepath); os.IsExist(node) {
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				fmt.Println("Things are not great")
				return
			}
			var p *types.PackageJSON
			json.Unmarshal(b, p)
			data.Projects[sub.Config().Name] = p
		}
	})
	return data, nil
}

// ResolveLocalDependancies will explore (via some graph expansion)
// once it is completed, it will return the list of the required
// pacakages otherwise, return an informative error
func (meta *Information) ResolveLocalDependancies(project string) ([]string, error) {
	if _, exists := meta.Projects[project]; !exists {
		return nil, errors.New("The project does not exists")
	}
	return []string{}, nil
}
