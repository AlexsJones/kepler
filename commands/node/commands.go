package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/AlexsJones/kepler/commands/submodules"

	git "gopkg.in/src-d/go-git.v4"
)

func recursePackages(p *PackageJSON, callback func(moduleName string, key string, value string)) error {

	for key, value := range p.Dependencies {

		callback(p.Name, key, value)
	}
	return nil
}

//HasPackage searches for packages references
//This is a useful way of whether a repository uses a packages
//It only requires package name without version
//On success it returns a bool and nil error object
func hasPackage(subPath string, filename string, target string) (bool, error) {

	filepath := path.Join(subPath, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false, err
	}
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return false, err
	}
	var packagejson PackageJSON
	json.Unmarshal(b, &packagejson)

	var wasFound = false
	recursePackages(&packagejson, func(moduleName string, key string, value string) {
		if strings.Contains(key, target) || strings.Contains(value, target) {
			wasFound = true
			fmt.Printf("Found usage in: %s, version is %s\n", moduleName, value)
			return
		}
	})
	return wasFound, nil
}

//FixLinks will perform a regex like action within a package.json to alter the url or file path
//It returns a nil error object on success
func fixLinks(subPath string, filename string, prefix string, target string, shouldDelete bool) error {

	filepath := path.Join(subPath, filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return err
	}
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	var packagejson PackageJSON
	json.Unmarshal(b, &packagejson)

	//processing
	recursePackages(&packagejson, func(moduleName string, key string, value string) {
		if strings.Contains(value, target) {
			if shouldDelete {
				delete(packagejson.Dependencies, key)
			} else {
				spli := strings.Split(value, "/")
				subspli := strings.Split(spli[len(spli[1:])], ".")
				foundEntry := subspli[0]
				foundEntry = strings.TrimSuffix(foundEntry, "\"")
				syntax := "file:%s%s"
				value := fmt.Sprintf(syntax, prefix, foundEntry)
				packagejson.Dependencies[key] = value
			}
		}
	})
	o, err := json.MarshalIndent(packagejson, "", "    ")
	if err != nil {
		return err
	}
	o = append(o, []byte("\n")...)

	return ioutil.WriteFile(filepath, o, 0644)
}

// NewInformation creates a struct containing information about the meta repo
func LocalNodeModules() (map[string]*PackageJSON, error) {
	Projects := make(map[string]*PackageJSON)
	submodules.LoopSubmodules(func(sub *git.Submodule) {
		filepath := path.Join(sub.Config().Path, "package.json")
		if _, node := os.Stat(filepath); !os.IsNotExist(node) {
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				return
			}
			var p PackageJSON
			json.Unmarshal(b, &p)
			Projects[sub.Config().Name] = &p
		}
	})
	return Projects, nil
}

// ResolveLocalDependancies will explore (via some graph expansion)
// once it is completed, it will return the list of the required
// pacakages otherwise, return an informative error
func ResolveLocalDependancies(project string) ([]string, error) {
	LocalPackages, err := LocalNodeModules()
	if err != nil {
		return []string{}, err
	}
	if _, exist := LocalPackages[project]; !exist {
		return nil, fmt.Errorf("%s does not exists", project)
	}
	ResolvedDeps := make(map[string]bool)
	Explore := make(map[string]*PackageJSON)
	// Making sure we don't try to explore the started node
	// if it is required by another project
	ResolvedDeps[project] = true
	Explore[project] = LocalPackages[project]
	for len(Explore) > 0 {
		for node, pack := range Explore {
			for name := range pack.Dependencies {
				if ResolvedDeps[name] {
					// Nothing to do as its already been resolved
					continue
				}
				if _, local := LocalPackages[name]; local {
					Explore[name] = LocalPackages[name]
				}
			}
			for name := range pack.DevDependencies {
				if ResolvedDeps[name] {
					// Nothing to do as its already been resolved
					continue
				}
				if _, local := LocalPackages[name]; local {
					Explore[name] = LocalPackages[name]
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
