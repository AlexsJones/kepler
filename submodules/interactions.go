package submodules

import (
	"fmt"
	"os"

	"github.com/AlexsJones/kepler/util"
	git "gopkg.in/libgit2/git2go.v25"
)

func loopSubmodules(path string, callback func(sub *git.Submodule, name string) error) error {
	repo, err := git.OpenRepository(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return repo.Submodules.Foreach(func(sub *git.Submodule, name string) int {

		callback(sub, name)
		return 0
	})
}

//UpdateSubmodules ...
func UpdateSubmodules(path string) error {

	loopSubmodules(path, func(sub *git.Submodule, name string) error {
		sub.Init(true)
		sub.Update(true, &git.SubmoduleUpdateOptions{
			CheckoutOpts: &git.CheckoutOpts{
				Strategy: git.CheckoutForce | git.CheckoutUpdateSubmodules,
			},
			FetchOptions:          &git.FetchOptions{},
			CloneCheckoutStrategy: git.CheckoutForce | git.CheckoutUpdateSubmodules | git.CheckoutSafe,
		})
		return nil
	})
	return nil
}

//CommandSubmodules ...
func CommandSubmodules(output string) error {

	loopSubmodules(".", func(sub *git.Submodule, name string) error {

		util.ShellCommand(output, sub.Path())

		return nil
	})
	return nil
}
