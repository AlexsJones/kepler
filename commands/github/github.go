package github

import (
	"context"
	"fmt"
	"log"

	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client
var localStorage *storage.Storage

//AddCommands for this module
func AddCommands(shell *ishell.Shell) string {

	shell.AddCmd(&ishell.Cmd{
		Name: "github-login",
		Help: "Login to github",
		Func: func(c *ishell.Context) {

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
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)
				c.Print("Access token: ")
				token := c.ReadPassword()
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
	})

	return "github"
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
