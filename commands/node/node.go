package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/submodules"
	"gopkg.in/src-d/go-git.v4"
)

//AddCommands for this module

//AddCommands for this module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "npm",
		Help: "npm command palette",
		Func: func(args []string) {
			fmt.Println("See help for working with npm")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "file",
				Help: "relink an npm package locally<prefix> <string>",
				Func: func(args []string) {
					if len(args) < 2 {
						fmt.Println("Please give a target package string to try to convert to a file link <prefix> <string> e.g. file ../../ googleremotes.git")
						return
					}
					submodules.LoopSubmodules(func(sub *git.Submodule) {
						if err := FixLinks(sub.Config().Path, "package.json", args[0], args[1], false); err != nil {
							fmt.Println(err.Error())
						} else {
							fmt.Printf("- Link fixed: %s\n", sub.Config().Path)
						}
					})
				},
			},
			command.Command{
				Name: "remove",
				Help: "remove a dep from package.json <string>",
				Func: func(args []string) {
					if len(args) < 1 {
						fmt.Println("Please give a target package string to to remove <string>")
						return
					}
					submodules.LoopSubmodules(func(sub *git.Submodule) {
						if err := FixLinks(sub.Config().Path, "package.json", "", args[0], true); err != nil {
						} else {
							fmt.Printf("- Removed in: %s\n", sub.Config().Path)
						}
					})
				},
			},
			command.Command{
				Name: "usage",
				Help: "find usage of a package within submodules <string>",
				Func: func(args []string) {
					if len(args) < 1 {
						fmt.Println("Find a package usage in submodule package.json <string> e.g. usage mocha")
						return
					}
					submodules.LoopSubmodules(func(sub *git.Submodule) {
						if _, err := HasPackage(sub.Config().Path, "package.json", args[0]); err != nil {
						}
					})
				},
			},
		},
	})

}

//PackageJSON structure of package.json
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Author          string            `json:"author"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func recursePackages(p *PackageJSON, callback func(moduleName string, key string, value string)) error {

	for key, value := range p.Dependencies {

		callback(p.Name, key, value)
	}
	return nil
}

//HasPackage searches for packages references
func HasPackage(subPath string, filename string, target string) (bool, error) {

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

//FixLinks ...
func FixLinks(subPath string, filename string, prefix string, target string, shouldDelete bool) error {

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
