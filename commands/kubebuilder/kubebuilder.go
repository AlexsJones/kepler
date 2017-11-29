package kubebuilder

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	event "github.com/AlexsJones/cloud-transponder/events"
	"github.com/AlexsJones/cloud-transponder/events/pubsub"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/AlexsJones/kubebuilder/src/data"
	"github.com/fatih/color"
	"github.com/gogo/protobuf/proto"
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

							fmt.Print("Please provide project name (e.g. my-gcloud-project):")
							reader := bufio.NewReader(os.Stdin)
							token, _ := reader.ReadString('\n')
							storage.GetInstance().Kubebuilder.ProjectName = strings.TrimSpace(token)

							fmt.Print("Please provide pubsub topic (e.g.cadium):")
							reader = bufio.NewReader(os.Stdin)
							token, _ = reader.ReadString('\n')
							storage.GetInstance().Kubebuilder.TopicName = strings.TrimSpace(token)

							fmt.Print("Please provide pubsub subscription (e.g.cadium-sub):")
							reader = bufio.NewReader(os.Stdin)
							token, _ = reader.ReadString('\n')
							storage.GetInstance().Kubebuilder.SubName = strings.TrimSpace(token)

							storage.GetInstance().Save()
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

func publishKubebuilderfile(build *data.BuildDefinition) error {

	//Create our GCP pubsub
	gpubsub := gcloud.NewPubSub()

	//Create the GCP Pubsub configuration
	gconfig := gcloud.NewPubSubConfiguration()

	gconfig.Topic = localStorage.Kubebuilder.TopicName
	gconfig.ConnectionString = localStorage.Kubebuilder.ProjectName
	gconfig.SubscriptionString = localStorage.Kubebuilder.SubName
	if err := event.Connect(gpubsub, gconfig); err != nil {
		return err
	}

	//Generate a new state object
	st := data.NewMessage(data.NewMessageContext())
	//Set our outbound message to indicate a build
	st.Type = data.Message_BUILD

	//Add the build as an encoded string into our message
	out, err := yaml.Marshal(build)
	if err != nil {
		return fmt.Errorf("Failed to marshal:%s", err)
	}

	st.Payload = base64.StdEncoding.EncodeToString(out)

	out, err = proto.Marshal(st)
	if err != nil {
		return fmt.Errorf("Failed to encode:%s", err)
	}

	err = event.Publish(gpubsub, out)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 5)

	color.Blue("Published to topic!")
	return nil
}
