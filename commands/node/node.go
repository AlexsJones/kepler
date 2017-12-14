//Package node provides a modular way of interacting with node commands
//This primarily is for dealing with nodejs files such as the package.json
package node

import (
	"fmt"
	"os"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	sh "github.com/AlexsJones/kepler/commands/shell"
	"github.com/AlexsJones/kepler/commands/submodules"
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
						deps, err := Resolve(project)
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
					defer func() {
						color.Yellow("Restoring backups")
						RestoreBackups()
					}()
					color.Yellow("Attempting to link packages")
					if err := LinkLocalDeps(); err != nil {
						color.Red("Failed to link: %s", err.Error())
						return
					}
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
					// Have to ensure that remove the old package.json
					// Otherwise there could be issues.
					os.Remove(filepath)
					if err = pack.WriteTo(filepath); err != nil {
						color.Red("Failed to write linked %s", filepath)
						color.Red("Due to %v", err)
					}
				},
			},
		},
	})
}
