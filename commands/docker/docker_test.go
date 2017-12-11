package docker_test

import (
	"testing"

	"github.com/AlexsJones/kepler/commands/docker"
)

func TestStandaloneCreation(t *testing.T) {
	template := `
FROM scratch

COPY {{.Application}} /src
{{range .Resources}}
SHOULD NOT EXISTS {{.}}
{{end}}`
	ExpectedTemplate := `
FROM scratch

COPY . /src
`
	config := &docker.Config{
		Application: "Test Application",
		Type:        "NoResolution",
		Resources:   []string{},
		Template:    []byte(template),
	}
	dockerfile, err := config.CreateStandaloneFile()
	if err != nil {
		t.Error("The template is rendered incorrectly: %v", err)
	}
	if string(dockerfile) != ExpectedTemplate {
		t.Log("Expected template:\n", ExpectedTemplate)
		t.Log("Given Template:\n", string(dockerfile))
		t.Error("Template was wrong!")
	}
}

func TestBadTemplate(t *testing.T) {
	template := `
FROM {{.Undefined}}

Plz break
	`
	config := &docker.Config{
		Application: "Test Application",
		Type:        "NoResolution",
		Resources:   []string{},
		Template:    []byte(template),
	}

	if _, err := config.CreateStandaloneFile(); err == nil {
		t.Error("Undefined template should report an error")
	}
}

func TestTemplateRequiringResources(t *testing.T) {
	template := `
FROM scratch:latest

COPY {{.Application}} /src/{{.Application}}
{{range .Resources}}
RUN echo {{.}}
{{end}}`
	expectedResult := `
FROM scratch:latest

COPY Test Application /src/Test Application

RUN echo 1

RUN echo 2

RUN echo 3
`
	config := &docker.Config{
		Application: "Test Application",
		Type:        "NoResolution",
		Resources:   []string{"1", "2", "3"},
		Template:    []byte(template),
	}
	result, err := config.CreateMetaFile()
	if err != nil {
		t.Error("Template file should not fail")
	}
	if string(result) != expectedResult {
		t.Log("Expected template:\n", expectedResult)
		t.Log("Given Template:\n", string(result))
		t.Error("Failed to create expected file")
	}
}

func TestUndefinedValues(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Code had to recover instead of fail gracefully")
		}
	}()
	config := &docker.Config{}
	if _, err := config.CreateStandaloneFile(); err == nil {
		t.Error("Failed to report on missing attributes")
	}
}
