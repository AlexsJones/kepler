package palette

import (
	"fmt"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/github"
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
						fmt.Println("Please login first...")
						return
					}
				},
			},
		},
	})
}
