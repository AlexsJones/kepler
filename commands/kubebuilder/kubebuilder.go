package kubebuilder

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/AlexsJones/kubebuilder/src/data"
	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v2"
)

var localStorage *storage.Storage

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
				Name: "setup",
				Help: "Configure the initial settings for kubebuilder",
				Func: func(args []string) {
					b, err := storage.Exists()
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					if b {
						//Load and save
						localStorage, err = storage.Load()
						if err != nil {
							color.Red(err.Error())
							return
						}
					} else {
						fmt.Print("Please install gcloud and authenticate (gcloud auth login) [Y/N] to continue:")
						reader := bufio.NewReader(os.Stdin)
						token, _ := reader.ReadString('\n')
						if strings.TrimSpace(token) == "Y" {
							localStorage = storage.NewStorage()
							fmt.Print("Please provide project name (e.g. my-gcloud-project):")
							reader := bufio.NewReader(os.Stdin)
							token, _ := reader.ReadString('\n')
							localStorage.Kubebuilder.ProjectName = strings.TrimSpace(token)

							fmt.Print("Please provide pubsub topic (e.g.cadium):")
							reader = bufio.NewReader(os.Stdin)
							token, _ = reader.ReadString('\n')
							localStorage.Kubebuilder.TopicName = strings.TrimSpace(token)

							storage.Save(localStorage)
						}
					}
					color.Green("Okay")
				},
			}, command.Command{
				Name: "deploy",
				Help: "Deploy to a remote kubebuilder cluster",
				Func: func(args []string) {
					if localStorage == nil {
						fmt.Println("Please run the setup first...")
						return
					}
					//--
					if _, err := loadKubebuilderFile(); err != nil {
						color.Red(err.Error())
					}
					color.Green("Okay")
				},
			},
		},
	},
	)
}

func loadKubebuilderFile() (*data.BuildDefinition, error) {

	if _, err := os.Stat(".kubebuilder"); os.IsNotExist(err) {
		return nil, errors.New(".kubebuilder folder does not exist")
	}
	if _, err := os.Stat(".kubebuilder/build.yaml"); os.IsNotExist(err) {
		return nil, errors.New(".kubebuilder folder does not exist")
	}

	//Load yaml
	raw, err := ioutil.ReadFile(".kubebuilder/build.yaml")
	if err != nil {
		log.Fatal(err)
	}
	//Hand cranking a build definition for the test
	builddef := data.BuildDefinition{}

	err = yaml.Unmarshal(raw, &builddef)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("%v\n", builddef)

	return &builddef, nil
}
