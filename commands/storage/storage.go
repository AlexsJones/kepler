package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	p "path"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

//AddCommands for this module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "storage",
		Help: "shell command palette",
		Func: func(args []string) {

		},
		SubCommands: []command.Command{
			command.Command{
				Name: "clear",
				Help: "clear all data from kepler",
				Func: func(args []string) {
					Delete()
					color.Blue("Deleted local storage")
				},
			},
		},
	})
}

var store = ".kepler"

//Storage structure
type Storage struct {
	Github *Github `json:"github"`
}

//Github specific sub structure
type Github struct {
	AccessToken string `json:"accesstoken"`
	IssueURL    string `json:"issueurl"`
}

//NewStorage object
func NewStorage() *Storage {

	s := &Storage{}
	s.Github = &Github{}
	return s
}

func path() (string, error) {
	s, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return p.Join(s, store), nil
}

//Exists in kepler
func Exists() (bool, error) {
	pout, err := path()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(pout); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

//Save to kepler storage
func Save(s *Storage) error {
	o, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}
	o = append(o, []byte("\n")...)
	pout, err := path()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pout, o, 0644)
}

//Delete storage
func Delete() error {
	pout, err := path()
	if err != nil {
		return err
	}
	b, _ := Exists()
	if b {
		os.Remove(pout)
	}
	return nil
}

//Load from kepler storage
func Load() (*Storage, error) {
	pout, err := path()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(pout)
	if err != nil {
		return nil, err
	}
	var s Storage
	json.Unmarshal(b, &s)

	return &s, nil
}
