package meta

import (
	"fmt"

	"github.com/Alexsjones/cli/cli"
	"github.com/Alexsjones/cli/command"
	"github.com/fatih/color"
)

//AddCommands for the submodule module
func AddCommands(cli *cli.Cli) {
	cli.AddCommand(command.Command{
		Name: "meta",
		Help: "Meta repo magic coming up",
		Func: func(args []string) {
			fmt.Println("See help for working with meta")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "node",
				Help: "observe all the wonderful commands we have for you today",
				Func: func(args []string) {
					if len(args) < 1 {
						fmt.Println("See help for working with branch")
						return
					}
				},
				SubCommands: []command.Command{
					command.Command{
						Name: "view",
						Help: "View all the node projects inside the node repo",
						Func: func(args []string) {
							i, err := NewInformation()
							if err != nil {
								color.Red("Something bad happened: %s", err.Error())
								return
							}
							for name := range i.Projects {
								color.Blue("> %s", name)
							}
						},
					},
				},
			},
		},
	})
}
