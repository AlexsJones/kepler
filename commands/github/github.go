package github

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client
var localStorage *storage.Storage

//AddCommands for this module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "github",
		Help: "github command palette",
		SubCommands: []command.Command{
			command.Command{
				Name: "login",
				Help: "use an access token to login to github",
				Func: func(args []string) {

					b, err := storage.Exists()
					if err != nil {
						fmt.Println(err.Error())
					}
					if b {
						//Load and save
						localStorage, err = storage.Load()
						if err != nil {
							return
						}
						log.Println("Loaded from storage")

					} else {
						fmt.Print("Access token: ")
						reader := bufio.NewReader(os.Stdin)
						token, _ := reader.ReadString('\n')
						log.Println("Creating new storage object...")
						localStorage = storage.NewStorage()
						localStorage.Github.AccessToken = token
						storage.Save(localStorage)
					}

					ctx := context.Background()
					ts := oauth2.StaticTokenSource(
						&oauth2.Token{AccessToken: localStorage.Github.AccessToken},
					)
					tc := oauth2.NewClient(ctx, ts)
					githubClient = github.NewClient(tc)
					_, _, err = githubClient.Repositories.List(ctx, "", nil)
					if err != nil {
						color.Red("Could not authenticate; please purge and login again")
						color.Red(err.Error())
						return
					}
					color.Green("Authentication Successful.")
				},
			},
		},
	})
}

//UnsetIssue from storage
func UnsetIssue() {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return
		}

	}
	localStorage.Github.IssueNumber = ""
	localStorage.Github.IssueRepo = ""
	storage.Save(localStorage)
}

//SetIssue in storage
func SetIssue(repo string, id string, store *storage.Storage) {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return
		}
	}
	store.Github.IssueNumber = id
	store.Github.IssueRepo = repo
	storage.Save(localStorage)
}

//AttachPRToIssue ...
func AttachPRToIssue() {
	//githubClient.PullRequests.CreateComment(ctx, owner, repo, number, i.GetURL())
}
