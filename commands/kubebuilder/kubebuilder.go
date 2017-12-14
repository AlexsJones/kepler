package kubebuilder

import (
	"fmt"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
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
					if err := Authenticate(); err != nil {
						color.Red("%v", err)
					} else {
						color.Green("Successfully logged in")
					}
				},
			},
			command.Command{
				Name: "build",
				Help: "Builds a docker image based off a kepler definitions",
				Func: func(args []string) {
					if len(args) == 0 {
						color.Red("Please type what projects you expect to build")
						return
					}
					for _, project := range args {
						if err := BuildDockerImage(project); err != nil {
							color.Red("%v", err)
							color.Yellow("If this is an auth issue, please make sure you have authenticated with gcloud")
							return
						}
					}
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
