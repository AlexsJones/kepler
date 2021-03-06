package docker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/AlexsJones/kepler/commands/node"
	sh "github.com/AlexsJones/kepler/commands/shell"
	yaml "gopkg.in/yaml.v2"
)

// Resolvers is a map of functions that will determine all
// the required external resources that could be found inside the meta repo
var Resolvers map[string]func(string) ([]string, error)

func init() {
	Resolvers = map[string]func(string) ([]string, error){
		// If we are going to build a node project from the meta
		// repo, we must enforce that we link all resolved projects.
		// Doing it inline as it shouldn't be called by any other
		// part of the application bare this part.
		"node": func(project string) ([]string, error) {
			node.LinkLocalDeps()
			return node.Resolve(project)
		},
	}
}

// Config contains all the required information
// in order to build the given application as a
// docker image
type Config struct {
	// Application is always assumed to be the base name of the current directory
	Application string
	// Type allows for correct resolution of required resources
	Type      string   `yaml:"Type"`
	BuildArgs []string `yaml:"BuildArgs"`
	Resources []string `yaml:"Resources"`
	Template  []byte
}

// CreateConfig loads the config defined in `ProjectDir/.kepler/config.yaml`
// and prepares the Dockerfile template defined in `ProjectDir/.kepler/Dockerfile`
// On success, it will return a struct with all the required information
// Otherwise, review the returned error message
func CreateConfig(projectDir string) (*Config, error) {
	// If we have been told to use the local directory,
	// Set it to blank so that we don't mess up references
	if projectDir == "." {
		projectDir = ""
	}
	conf := path.Join(projectDir, ".kepler/config.yaml")
	if _, err := os.Stat(conf); os.IsNotExist(err) {
		return nil, fmt.Errorf("Unable to find %s", conf)
	}
	b, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	config := Config{
		Application: path.Base(projectDir),
	}
	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	// Enforce that resources are not resolved for this if
	// config type isn't defined
	template := path.Join(projectDir, ".kepler/Dockerfile.tmpl")
	if _, err = os.Stat(template); os.IsNotExist(err) {
		return nil, fmt.Errorf("Expected file %s missing", template)
	}
	b, err = ioutil.ReadFile(template)
	if err != nil {
		return nil, err
	}
	config.Template = b
	return &config, nil
}

func (conf *Config) prepareTemplate() ([]byte, error) {
	resolverType := strings.ToLower(conf.Type)
	resolverFunc, exist := Resolvers[resolverType]
	if exist {
		addtional, err := resolverFunc(conf.Application)
		if err != nil {
			return nil, err
		}
		conf.Resources = append(conf.Resources, addtional...)
	}
	t := template.Must(template.New("Dockerfile").Parse(string(conf.Template)))
	dockerfile := &bytes.Buffer{}
	err := t.Execute(dockerfile, conf)
	return dockerfile.Bytes(), err
}

// CreateStandaloneFile strips all the templating from the original template
// without doing any resource resolution and returns a byte stream that is used as the dockerfile
func (conf *Config) CreateStandaloneFile() ([]byte, error) {
	// Have to modify the application name to be dot to
	// ensure it copies the files in the current directory
	if err := conf.validate(); err != nil {
		return nil, err
	}
	name := conf.Application
	conf.Application = "."
	conf.Type = ""
	b, err := conf.prepareTemplate()
	conf.Application = name
	return b, err
}

// CreateMetaFile will create a dockerfile based off the config.Type
// and the resources it requires
func (conf *Config) CreateMetaFile() ([]byte, error) {
	if err := conf.validate(); err != nil {
		return nil, err
	}
	return conf.prepareTemplate()
}

func BuildImage(buildArgs ...string) error {
	// Once someone creates an actual working golang api to interact with Docker
	// this is the best I can do.
	return sh.ShellCommand(fmt.Sprintf("docker build %s .", strings.Join(buildArgs, " ")), ".", false)
}

func (conf *Config) validate() error {
	if conf.Application == "" {
		return fmt.Errorf("Application does not have a valid value")
	}
	return nil
}
