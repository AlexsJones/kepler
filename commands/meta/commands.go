package meta

import (
	"fmt"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/fatih/color"
)

//AddCommands for the Meta module
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
						fmt.Println("See help for working with meta node")
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
							if len(i.Projects) == 0 {
								color.Red("There appears to be no happiness in the world")
								return
							}
							for name := range i.Projects {
								color.Blue("> %s", name)
							}
						},
					},
					command.Command{
						Name: "local-deps",
						Help: "Shows all the dependancies found locally",
						Func: func(args []string) {
							i, err := NewInformation()
							if err != nil {
								color.Red("Something bad has happened: %s", err.Error())
								return
							}
							for _, project := range args {
								deps, err := i.ResolveLocalDependancies(project)
								if err != nil {
									color.Red("The hell?!: %s", err.Error())
								} else {
									color.Cyan("> %s", project)
									for _, dep := range deps {
										color.Green("> %s", dep)
									}
								}

							}
						},
					},
				},
			},
		},
	})
}
