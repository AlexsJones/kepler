package submodules

import (
	"fmt"
	"os"
	"os/exec"

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
		sub.FetchRecurseSubmodules()

		return nil
	})
	return nil
}

//CommandSubmodules ...
func CommandSubmodules(path string, output string) error {

	loopSubmodules(path, func(sub *git.Submodule, name string) error {
		cmd := exec.Command("bash", "-c", output)
		cmd.Dir = sub.Path()
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf(string(out))
		return nil
	})

	return nil
}
