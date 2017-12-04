//Package storage is a file system abstraction for storing Kepler data
package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	p "path"
	"sync"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

//AddCommands for this module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "storage",
		Help: "storage command palette",
		Func: func(args []string) {
			fmt.Println("See help for working with shell")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "clear",
				Help: "clear all data from kepler",
				Func: func(args []string) {
					Delete()
					color.Green("Deleted local storage")
				},
			},
			command.Command{
				Name: "show",
				Help: "show storage data",
				Func: func(args []string) {
					ShowStorage()
				},
			},
		},
	})
}

var store = ".kepler"

//Storage structure
type Storage struct {
	Github      *Github      `json:"github"`
	Kubebuilder *Kubebuilder `json:"kubebuilder"`
}

//Kubebuilder specific sub structure
type Kubebuilder struct {
	ProjectName string `json:"projectname"`
	TopicName   string `json:"topicname"`
	SubName     string `json:"subscriptionname"`
}

//Github specific sub structure
type Github struct {
	AccessToken  string  `json:"accesstoken"`
	Issue        []Issue `json:"issue"`
	TeamID       int     `json:"teamid"`
	CurrentIssue *Issue  `json:"currentissue"`
}

//Issue object structure
type Issue struct {
	IssueURL     string        `json:"issueurl"`
	Owner        string        `json:"owner"`
	Repo         string        `json:"repo"`
	Number       int           `json:"number"`
	PullRequests []PullRequest `json:"pullrequests"`
	Palette      map[string]string
}

//PullRequest object structure
type PullRequest struct {
	Repo   string
	Owner  string
	Base   string
	Head   string
	Title  string
	Number int
}

var instance *Storage
var once sync.Once

//GetInstance reference to the singleton
func GetInstance() *Storage {
	once.Do(func() {
		doesExist, err := Exists()
		if err != nil {
			panic(err)
		}
		if doesExist != true {
			instance = &Storage{}
			instance.Github = &Github{TeamID: 0}
			instance.Kubebuilder = &Kubebuilder{}
			log.Println("Creating new storage object...")
		} else {
			i, err := Load()
			if err != nil {
				panic(err)
			}
			instance = i
			log.Println("Using existing storage object...")
		}
	})
	return instance
}

func path() (string, error) {
	s, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return p.Join(s, store), nil
}

//Exists checks if .kepler file has been set
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
func (s *Storage) Save() error {
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

//ShowStorage in kepler
func ShowStorage() error {
	pout, err := path()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(pout)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
