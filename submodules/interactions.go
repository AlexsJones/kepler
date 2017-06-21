package submodules

import (
	"fmt"
	"os"

	"github.com/AlexsJones/kepler/util"

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
	return nil
}

//CommandSubmodules ...
func CommandSubmodules(output string) error {

	loopSubmodules(".", func(sub *git.Submodule) error {

		util.ShellCommand(output, sub.Config().Path)

		return nil
	})
	return nil
}
