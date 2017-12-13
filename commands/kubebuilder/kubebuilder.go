package kubebuilder

import (
	"context"
	"flag"
	"fmt"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	gcr "github.com/GoogleCloudPlatform/docker-credential-gcr/cli"
	"github.com/fatih/color"
)

//AddCommands for the kubebuilder module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "kubebuilder",
		Help: "kubebuilder command palette",
		Func: func(args []string) {
			fmt.Println("See help for working with kubebuilder")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "auth",
				Help: "Authenticates you against all required services",
				Func: func(args []string) {
					gcr.NewGCRLoginSubcommand().Execute(context.Background(), &flag.FlagSet{})
				},
			},
			command.Command{
				Name: "deploy",
				Help: "Deploy to a remote kubebuilder cluster",
				Func: func(args []string) {

					out, err := loadKubebuilderFile()
					if err != nil {
						color.Red(err.Error())
						return
					}

					if err := publishKubebuilderfile(out); err != nil {
						color.Red(err.Error())
						return
					}

					color.Green("Okay")
				},
			},
		},
	},
	)
}
