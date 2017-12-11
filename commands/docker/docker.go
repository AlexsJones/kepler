package docker

import (
	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/fatih/color"
)

func AddCommands(cli *cli.Cli) {
	cli.AddCommand(command.Command{
		Name: "docker",
		Help: "docker command palette",
		Func: func(args []string) {
			color.Magenta("WIP")
		},
	})
}
