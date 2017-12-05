// Package meta is aimed at working
// within a git repo that manages other repos.
// The issues this package helps solve is that
// node packages are a pain with npm as you
// can't say use local if you can find them
// instead of using remote.
package meta

import (
	"errors"
	"os"
)

// Information ...
type Information struct {
	// projects, poor mans version of a set
	projects map[string]bool
}

// NewInformation creates a struct containing information about the meta repo
func NewInformation() (*Information, error) {
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		return nil, errors.New("Can not create information about meta repo")
	}
	data := &Information{}
	return data, nil
}

// ResolveLocalDependancies will explore (via some graph expansion)
// once it is completed, it will return the list of the required
// pacakages otherwise, return an informative error
func (meta *Information) ResolveLocalDependancies(project string) ([]string, error) {
	if !meta.projects[project] {
		return nil, errors.New("The project does not exists")
	}
	return []string{}, nil
}
