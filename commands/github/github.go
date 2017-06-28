package github

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client
var ctx context.Context
var localStorage *storage.Storage

//AddCommands for this module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "github",
		Help: "github command palette",
		Func: func(args []string) {
			fmt.Println("See help for working with github")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "pr",
				Help: "pr command palette",
				Func: func(args []string) {
					fmt.Println("See help for working with pr")
				},
				SubCommands: []command.Command{
					command.Command{
						Name: "attach",
						Help: "attach the current issue to a pr <reponame> <owner> <prnumber>",
						Func: func(args []string) {
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							if len(args) == 0 || len(args) < 3 {
								fmt.Println("set the current working issue in the pr <reponame> <owner> <prnumber>")
								return
							}
							AttachIssuetoPr(args[0], args[1], args[2])
						},
					},
				},
			},
			command.Command{
				Name: "issue",
				Help: "Issue command palette",
				Func: func(args []string) {
					fmt.Println("See help for working with issue")
				},
				SubCommands: []command.Command{
					command.Command{
						Name: "set",
						Help: "set the current working issue <issue url>",
						Func: func(args []string) {
							if len(args) == 0 || len(args) < 1 {
								fmt.Println("Requires <issue url>")
								return
							}
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							SetIssue(args[0])
							color.Green("Okay")
						},
					},
					command.Command{
						Name: "unset",
						Help: "unset the current working issue",
						Func: func(args []string) {
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							UnsetIssue()
						},
					},
				},
			},
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

					ctx = context.Background()
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
	localStorage.Github.IssueURL = ""
	storage.Save(localStorage)
}

//SetIssue in storage
func SetIssue(issueurl string) {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return
		}
	}
	localStorage.Github.IssueURL = issueurl
	storage.Save(localStorage)
}

//AttachIssuetoPr ...
func AttachIssuetoPr(reponame string, owner string, number string) {

	if localStorage == nil {
		localStorage, _ = storage.Load()
	}

	if localStorage.Github.IssueURL == "" {
		color.Red("No working issue set...")
		return
	}

	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(err)
	}

	pr, res, err := githubClient.PullRequests.Get(ctx, owner, reponame, num)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Github says %d\n", res.StatusCode)

	appended := fmt.Sprintf("%s\n%s\n", string(pr.GetBody()), localStorage.Github.IssueURL)

	pr, res, err = githubClient.PullRequests.Edit(ctx, owner, reponame, num, &github.PullRequest{Body: &appended})
	if err != nil {
		fmt.Println(err)
	}
	color.Green("Okay")
}
