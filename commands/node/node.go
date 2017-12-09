//Package node provides a modular way of interacting with node commands
//This primarily is for dealing with nodejs files such as the package.json
package node

import (
	"fmt"
	"io/ioutil"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	sh "github.com/AlexsJones/kepler/commands/shell"
	"github.com/AlexsJones/kepler/commands/submodules"
	"github.com/MovieStoreGuy/resources/marshal"
	"github.com/fatih/color"
	"gopkg.in/src-d/go-git.v4"
)

//AddCommands for the node module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "node",
		Help: "node command palette",
		Func: func(args []string) {
			fmt.Println("See help for working with npm")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "remove",
				Help: "remove a dep from package.json <string>",
				Func: func(args []string) {
					if len(args) < 1 {
						fmt.Println("Please give a target package string to to remove <string>")
						return
					}
					submodules.LoopSubmodules(func(sub *git.Submodule) {
						if err := fixLinks(sub.Config().Path, "package.json", "", args[0], true); err != nil {
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
						if _, err := hasPackage(sub.Config().Path, "package.json", args[0]); err != nil {
						}
					})
				},
			},
			command.Command{
				Name: "view",
				Help: "View all the node projects that can be found locally",
				Func: func(args []string) {
					i, err := LocalNodeModules()
					if err != nil {
						color.Red("Something bad happened: %s", err.Error())
						return
					}
					if len(i) == 0 {
						color.Red("No submodules found")
						return
					}
					for name := range i {
						color.Blue("> %s", name)
					}
				},
			},
			command.Command{
				Name: "local-deps",
				Help: "Shows all the dependancies found locally",
				Func: func(args []string) {
					for _, project := range args {
						deps, err := ResolveLocalDependancies(project)
						if err != nil {
							color.Red("The hell?!: %s", err.Error())
						} else {
							color.Cyan("> %s", project)
							for _, dep := range deps {
								color.Green("> %s", dep)
							}
						}

					}
				},
			},
			command.Command{
				Name: "install",
				Help: "Installs all the required vendor code",
				Func: func(args []string) {
					color.Yellow("Attempting to link packages")
					if err := LinkLocalDeps(); err != nil {
						color.Red("Failed to link: %s", err.Error())
						return
					}
					defer func() {
						color.Yellow("Restoring backups")
						RestoreBackups()
					}()
					color.Yellow("Attempting to install")
					sh.ShellCommand("npm i", "", true)
				},
			},
			command.Command{
				Name: "init",
				Help: "Create the package json for a meta repo",
				Func: func(args []string) {
					pack, err := CreateMetaPackageJson()
					if err != nil {
						color.Red("Failed to generate meta package json, %s", err.Error())
						return
					}
					// Write new package json to disk
					filepath := "package.json"
					o, err := marshal.PureMarshalIndent(pack, "", "    ")
					if err != nil {
						color.Red("An error occured, %s", err.Error())
						return
					}
					o = append(o, []byte("\n")...)
					if err = ioutil.WriteFile(filepath, o, 0644); err != nil {
						color.Red("Failed to write linked %s", filepath)
					}
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
	Bugs            map[string]string `json:"bugs,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
	Private         bool              `json:"private,omitempty"`
	License         string            `json:"license,omitempty"`
}
