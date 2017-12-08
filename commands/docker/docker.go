package docker

import (
	"io/ioutil"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/fatih/color"
)

func AddCommands(cli *cli.Cli) {
	cli.AddCommand(command.Command{
		Name: "docker",
		Help: "docker command palette",
		Func: func(args []string) {},
		SubCommands: []command.Command{
			command.Command{
				Name: "create",
				Help: "Creates a dockerfile from a template",
				Func: func(args []string) {
					if len(args) == 0 {
						color.Red("Requires an argument to know what to build")
						return
					}
					project := args[0]
					dockerfile, err := CreateDockerfile(project)
					if err != nil {
						color.Red("An issue occured: %s", err.Error())
						return
					}
					if err = ioutil.WriteFile("Dockerfile", dockerfile, 0644); err != nil {
						color.Red("%s", err.Error())
					} else {
						color.Green("Successfully created Dockerfile for %s", project)
					}
				},
			},
		},
	})
}
