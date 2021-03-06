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
	"github.com/AlexsJones/kepler/commands/docker"
	"github.com/AlexsJones/kepler/commands/node"
	sh "github.com/AlexsJones/kepler/commands/shell"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/AlexsJones/kubebuilder/src/data"
	login "github.com/GoogleCloudPlatform/docker-credential-gcr/auth"
	"github.com/GoogleCloudPlatform/docker-credential-gcr/config"
	"github.com/fatih/color"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/oauth2"
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
		return nil, err
	}
	//Hand cranking a build definition for the test
	builddef := data.BuildDefinition{}

	err = yaml.Unmarshal(raw, &builddef)
	if err != nil {
		return nil, err
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

func authenticateDocker() error {
	for _, registry := range []string{"gcr.io", "us.gcr.io"} {
		command := fmt.Sprintf("docker login -u %s -p %s https://%s", "oauth2accesstoken", storage.GetInstance().GCRAuth.AccessToken, registry)
		// Gross hack untill "github.com/moby/moby" can be fetched using go get
		if err := sh.ShellCommand(command, "", false); err != nil {
			return err
		}
	}
	return nil
}

// BuildDockerImage will load the config within the given directory
// and will build an image based on those parameters
func BuildDockerImage(project string) error {
	// If a Dockerfile lives in the current directory,
	// we can not assume that it has all the current information so we have to
	// abort and let the callee resolve this issue.
	if _, err := os.Stat("Dockerfile"); !os.IsNotExist(err) {
		return fmt.Errorf("Dockerfile found within local directory, aborting")
	}
	config, err := docker.CreateConfig(project)
	if err != nil {
		return err
	}

	// When resolving different config types, we may also be rewriting
	// content on disc, so this should ensure that content as the user left it.
	switch config.Type {
	case "node":
		defer node.RestoreBackups()
	}

	dockerfile, err := config.CreateMetaFile()
	if err != nil {
		return err
	}
	// We want to remove the generated Dockefiler once we are done
	defer os.Remove("Dockerfile")
	if err := ioutil.WriteFile("Dockerfile", dockerfile, 0644); err != nil {
		return err
	}
	return nil
}

// Authenticate will login to the required services only if the services
// keys have expired or require updating.
func Authenticate() error {
	auth := storage.GetInstance().GCRAuth
	switch {
	case auth == nil:
		// Using the GCR Login Agent to obtain us
		// the required access token
		client := &login.GCRLoginAgent{
			AllowBrowser: true,
		}
		tkn, err := client.PerformLogin()
		if err != nil {
			return err
		}
		// Gross casting due to vendor code being used in a third party package
		var cast interface{} = tkn
		token, ok := cast.(*oauth2.Token)
		if !ok {
			return fmt.Errorf("Unable to correctly cast to *oauth.Token")
		}
		storage.GetInstance().GCRAuth = token
	case time.Now().After(auth.Expiry):
		color.Yellow("Auth has expired, trying to refresh token")
		var cast interface{} = config.GCROAuth2Endpoint
		endpoint, ok := cast.(oauth2.Endpoint)
		if !ok {
			return fmt.Errorf("Unable to oauth2.EndPoint")
		}
		conf := &oauth2.Config{
			ClientID:     config.GCRCredHelperClientID,
			ClientSecret: config.GCRCredHelperClientNotSoSecret,
			Scopes:       config.GCRScopes,
			Endpoint:     endpoint,
		}
		// It is expected that will update our access token instead of
		// constantly asking for us to update

		token, err := oauth2.ReuseTokenSource(auth, conf.TokenSource(config.OAuthHTTPContext, auth)).Token()
		if err != nil {
			return err
		}
		storage.GetInstance().GCRAuth = token
	}
	return storage.GetInstance().Save()
}
