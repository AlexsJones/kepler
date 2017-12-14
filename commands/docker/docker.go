package docker

import (
	"io/ioutil"
	"os"

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
		SubCommands: []command.Command{
			command.Command{
				Name: "build",
				Help: "Builds a project in standalone from the defined Dockerfile",
				Func: func(args []string) {
					if _, err := os.Stat("Dockerfile"); os.IsNotExist(err) {
						color.Blue("No Dockerfile found locally")
						color.Blue("Attempting to build a dockerfile from a template")
						config, err := CreateConfig(".")
						if err != nil {
							color.Red("%v", err)
							return
						}
						dockerfile, err := config.CreateStandaloneFile()
						if err != nil {
							color.Red("%v", err)
							return
						}
						if err = ioutil.WriteFile("Dockerfile", dockerfile, 0644); err != nil {
							color.Red("%v", err)
							return
						}
						// Make sure we remove our templated Dockerfile once we are done
						defer os.Remove("Dockerfile")
					}
					BuildImage(args...)
				},
			},
		},
	})
}
