package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"gopkg.in/src-d/go-git.v4"
)

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

		ShellCommand(output, sub.Config().Path, false)

		return nil
	})
	return nil
}
