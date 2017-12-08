package docker

import (
	"strings"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/node"
	"github.com/fatih/color"
)

func AddCommands(cli *cli.Cli) {
	cli.AddCommand(command.Command{
		Name: "docker",
		Help: "docker command palette",
		Func: func(args []string) {},
		SubCommands: []command.Command{
			command.Command{
				Name: "build",
				Help: "Builds a docker image from a template",
				Func: func(args []string) {
					if len(args) == 0 {
						color.Blue("build requires arguments")
						return
					}
					node.LinkLocalDeps()
					defer node.RestoreBackups()
					for _, project := range args {
						color.Green("Starting to build %s", project)
						if err := BuildImage(project); err != nil {
							color.Red("Problem occured: %s", err.Error())
							return
						}
					}
				},
			},
			command.Command{
				Name: "build-args",
				Help: "Sets the required build args for kepler",
				Func: func(args []string) {
					BuildArgs = strings.Join(args, " ")
					color.Green("Build args is now set to: %s", BuildArgs)
				},
			},
		},
	})
}
