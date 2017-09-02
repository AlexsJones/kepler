package palette

import (
	"fmt"
	"os"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/github"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/fatih/color"
)

//AddCommands for the palette module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "palette",
		Help: "Issue palette that controls repos can be used from here",
		Func: func(args []string) {
			fmt.Println("See help for working with palette")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "branch",
				Help: "switch branches or create if they don't exist for working issue palette repos <branchname>",
				Func: func(args []string) {
					if github.GithubClient == nil || github.LocalStorage == nil {
						//Warning: Not sure this works yet
						fmt.Println("Please login first...")
						return
					}
					if len(args) == 0 || len(args) < 1 {
						fmt.Println("provide the branch name to switch repo in the palette too <branchname>")
						return
					}
					for k, v := range github.LocalStorage.Github.CurrentIssue.Palette {
						if _, err := os.Stat(v); os.IsNotExist(err) {
							color.Red(fmt.Sprintf("Warning the repo %s does not exist at the path %s, removing from the palette\n", k, v))
							delete(github.LocalStorage.Github.CurrentIssue.Palette, k)
							storage.Save(github.LocalStorage)
						} else {
							color.Green(fmt.Sprintf("Switching %s to branch %s", k, args[0]))

						}
					}
				},
			},
		},
	})
}
