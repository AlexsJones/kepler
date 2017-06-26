package submodules

import (
	"fmt"
	"os"
	"strings"

	sh "github.com/AlexsJones/kepler/commands/shell"
	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
	"gopkg.in/src-d/go-git.v4"
)

//AddCommands to this module
func AddCommands(shell *ishell.Shell) string {
	shell.AddCmd(&ishell.Cmd{
		Name: "submodules-exec",
		Help: "Exec command in submodules <cmd> e.g. exec git reset --hard HEAD",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Please provide a command")
				return
			}
			CommandSubmodules(strings.Join(c.Args, " "))
		},
	})
	return "submodules"
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
