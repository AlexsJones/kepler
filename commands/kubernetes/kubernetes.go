package kubernetes

import (
	"fmt"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
)

//AddCommands for this module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "kubernetes",
		Help: "kubernetes command palette",
		Func: func(args []string) {
			fmt.Println("See help for working with kubernetes & kubectl")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "generate",
				Help: "generate templates to make your life easier",
				Func: func(args []string) {

				},
				SubCommands: []command.Command{
					command.Command{
						Name: "deployment",
						Help: "generate a barebones deployment with your deployment name <name>",
						Func: func(args []string) {
							if len(args) == 0 || len(args) < 1 {
								fmt.Println("Requires <name>")
								return
							}

						},
					},
				},
			},
		},
	})
}
