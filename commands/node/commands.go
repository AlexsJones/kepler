package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlexsJones/kepler/commands/submodules"
	"github.com/MovieStoreGuy/resources/files"
	"github.com/fatih/color"

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

// LocalNodeModules creates a struct containing information about the meta repo
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

func LinkLocalDeps() error {
	local, err := LocalNodeModules()
	if err != nil {
		return err
	}
	for dir, pack := range local {
		// Need to create a backup file
		filepath := path.Join(dir, "package.json.bak")
		if _, err = os.Stat(filepath); !os.IsNotExist(err) {
			color.Red("%s already exists", filepath)
			continue
		}
		if err = files.Copy(path.Join(dir, "package.json"), filepath); err != nil {
			color.Red("Failed to create %s", filepath)
		}
		color.Blue("Updating %s links", dir)
		for name := range pack.Dependencies {
			if _, exist := local[name]; exist {
				pack.Dependencies[name] = fmt.Sprintf("file:../%s", name)
			}
		}
		for name := range pack.DevDependencies {
			if _, exist := local[name]; exist {
				pack.DevDependencies[name] = fmt.Sprintf("file:../%s", name)
			}
		}
		// Write new package json to disk
		filepath = path.Join(dir, "package.json")
		if err := os.Remove(filepath); err != nil {
			return err
		}
		o, err := json.MarshalIndent(pack, "", "    ")
		if err != nil {
			return err
		}
		o = append(o, []byte("\n")...)
		if err = ioutil.WriteFile(filepath, o, 0644); err != nil {
			color.Red("Failed to write linked %s", filepath)
			return err
		}
	}
	return nil
}

func RestoreBackups() error {
	local, err := LocalNodeModules()
	if err != nil {
		return err
	}
	for name := range local {
		filepath := path.Join(name, "package.json.bak")
		if _, err = os.Stat(filepath); !os.IsNotExist(err) {
			if err = files.Copy(filepath, path.Join(name, "package.json")); err != nil {
				return err
			}
			// Need to remove the packup file
			if err = os.Remove(filepath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateMetaPackageJson() (*PackageJSON, error) {
	metaPackage := &PackageJSON{
		Version:         "1.0.0",
		Description:     "An auto generated package json",
		Main:            "index.js",
		Author:          os.Args[0],
		Dependencies:    map[string]string{},
		DevDependencies: map[string]string{},
		Scripts: map[string]string{
			"test": "true",
		},
	}
	if name, err := os.Getwd(); err != nil {
		return nil, err
	} else {
		metaPackage.Name = filepath.Base(name)
	}
	modules, err := LocalNodeModules()
	if err != nil {
		return nil, err
	}
	for name := range modules {
		metaPackage.Dependencies[name] = fmt.Sprintf("file:%s", name)
	}
	return metaPackage, nil
}
