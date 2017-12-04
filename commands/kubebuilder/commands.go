package kubebuilder

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	event "github.com/AlexsJones/cloud-transponder/events"
	"github.com/AlexsJones/cloud-transponder/events/pubsub"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/AlexsJones/kubebuilder/src/data"
	"github.com/fatih/color"
	"github.com/gogo/protobuf/proto"
	yaml "gopkg.in/yaml.v2"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
func loadKubebuilderFile() (*data.BuildDefinition, error) {

	if _, err := exists(".kubebuilder"); os.IsNotExist(err) {
		return nil, errors.New(".kubebuilder folder does not exist")
	}
	if _, err := exists(".kubebuilder/build.yaml"); os.IsNotExist(err) {
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

	gconfig.Topic = storage.GetInstance().Kubebuilder.TopicName
	gconfig.ConnectionString = storage.GetInstance().Kubebuilder.ProjectName
	gconfig.SubscriptionString = storage.GetInstance().Kubebuilder.SubName
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
