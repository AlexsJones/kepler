package kubebuilder

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/fatih/color"
	rawdocker "github.com/fsouza/go-dockerclient"
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
					access, err := authenticateDocker()
					if err != nil {
						color.Red("Failed to login")
						color.Red("Due to %v", err)
						return
					} else {
						color.Green("We are logged in to GCR")
					}
					client, _ := rawdocker.NewClientFromEnv()
					t := time.Now()
					dockerfile, err := ioutil.ReadFile("Dockerfile")
					if err != nil {
						color.Red("%v", err)
						return
					}
					dockerfile = append(dockerfile, []byte("\n")...)
					inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
					tr := tar.NewWriter(inputbuf)
					tr.WriteHeader(&tar.Header{
						Name:       "Dockerfile",
						Size:       int64(len(dockerfile) + 1),
						ModTime:    t,
						AccessTime: t,
						ChangeTime: t,
					})
					tr.Write(dockerfile)
					tr.Close()
					opts := rawdocker.BuildImageOptions{
						Name:         "api-core",
						InputStream:  inputbuf,
						OutputStream: outputbuf,
						Auth:         *access,
						AuthConfigs: rawdocker.AuthConfigurations{
							Configs: map[string]rawdocker.AuthConfiguration{
								"https://us.gcr.io": *access,
							},
						},
					}
					if err := client.BuildImage(opts); err != nil {
						color.Red("We gone fucked up")
						color.Red("%v", err)
					}
				},
			},
			command.Command{
				Name: "build",
				Help: "Builds a docker image based off a kepler definitions",
				Func: func(args []string) {
					if len(args) == 0 {
						color.Red("Please tpye what projects you expect to build")
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
