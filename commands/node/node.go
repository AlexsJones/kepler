//Package node provides a modular way of interacting with node commands
//This primarily is for dealing with nodejs files such as the package.json
package node

import (
	"fmt"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/submodules"
	"gopkg.in/src-d/go-git.v4"
)

//AddCommands for the node module
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
						if err := fixLinks(sub.Config().Path, "package.json", args[0], args[1], false); err != nil {
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
