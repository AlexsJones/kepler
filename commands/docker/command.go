package docker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/AlexsJones/kepler/commands/node"
)

type Resources struct {
	Application string
	Resources   []string
}

var TemplateDirectory string

func init() {
	TemplateDirectory = "templates"
}

func CreateDockerfile(application string) ([]byte, error) {
	file := path.Join(TemplateDirectory, fmt.Sprintf("%s.Dockerfile", application))
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, fmt.Errorf("Can not load %s", file)
	}
	resources := &Resources{
		Application: application,
	}
	if deps, err := node.ResolveLocalDependancies(application); err != nil {
		return nil, err
	} else {
		resources.Resources = deps
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	t := template.Must(template.New("Dockerfile").Parse(string(b)))
	dockerfile := &bytes.Buffer{}
	t.Execute(dockerfile, resources)
	return dockerfile.Bytes(), nil
}
