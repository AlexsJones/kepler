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
	yaml "gopkg.in/yaml.v2"
)

var Resolvers map[string]func(string) ([]string, error)

func init() {
	Resolvers = map[string]func(string) ([]string, error){
		"node": node.Resolve,
		"noresolution": func(empty string) ([]string, error) {
			return []string{}, nil
		},
	}
}

// Config contains all the required information
// in order to build the given application as a
// docker image
type Config struct {
	Application string
	// Type allows for correct resolution of required resources
	Type      string   `yaml:"Type"`
	BuildArgs []string `yaml:"BuildArgs"`
	Resources []string
	Template  []byte
}

// CreateConfig loads the config defined in `ProjectDir/.kepler/config.yaml`
// and prepares the Dockerfile template defined in `ProjectDir/.kepler/Dockerfile`
// On success, it will return a struct with all the required information
// Otherwise, review the returned error message
func CreateConfig(ProjectDir string) (*Config, error) {
	conf := path.Join(ProjectDir, ".kepler/config.yaml")
	if _, err := os.Stat(conf); os.IsNotExist(err) {
		return nil, fmt.Errorf("Unable to find %s", conf)
	}
	b, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	config := Config{
		Application: ProjectDir,
	}
	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	template := path.Join(ProjectDir, ".kepler/Dockerfile")
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
	if _, exist := Resolvers[resolverType]; !exist {
		return nil, fmt.Errorf("Undefined project type")
	}
	if resources, err := Resolvers[resolverType](conf.Application); err != nil {
		return nil, err
	} else {
		conf.Resources = resources
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
	conf.Type = "noresolution"
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

func (conf *Config) validate() error {
	if conf.Application == "" {
		return fmt.Errorf("Application does not have a valid value")
	}
	return nil
}
