package submodules

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	sh "github.com/AlexsJones/kepler/commands/shell"
	"github.com/fatih/color"
	"gopkg.in/src-d/go-git.v4"
)

//AddCommands for this module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "submodule",
		Help: "submodule command palette",
		Func: func(args []string) {
			fmt.Println("See help for working with submodules")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "branch",
				Help: "branch command palette",
				Func: func(args []string) {
					if len(args) < 1 {
						fmt.Println("See help for working with branch")
						return
					}
				},
			},
			command.Command{
				Name: "exec",
				Help: "execute in all submodules <command string>",
				Func: func(args []string) {
					if len(args) < 1 {
						fmt.Println("Please provide a command <command string>")
						return
					}
					CommandSubmodules(strings.Join(args, " "))
				},
			},
		},
	})
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

//LoopSubmodules ...
func LoopSubmodules(callback func(sub *git.Submodule)) error {
	loopSubmodules(".", func(sub *git.Submodule) error {
		callback(sub)
		return nil
	})
	return nil
}

//CommandSubmodules ...
func CommandSubmodules(output string) error {

	loopSubmodules(".", func(sub *git.Submodule) error {

		sh.ShellCommand(output, sub.Config().Path, false)

		return nil
	})
	return nil
}
