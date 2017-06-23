package github

import (
	"context"
	"fmt"
	"log"

	"github.com/AlexsJones/kepler/util"
	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client

//AddCommands for this module
func AddCommands(shell *ishell.Shell) {

	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Help: "Login to github",
		Func: func(c *ishell.Context) {

			var localStorage *storage.Storage
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
}
