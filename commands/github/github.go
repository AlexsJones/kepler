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
					command.Command{
						Name: "create",
						Help: "create a pr <reponame> <owner> <title> <base> <head>",
						Func: func(args []string) {
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							if len(args) == 0 || len(args) < 5 {
								fmt.Println("create a pr <reponame> <owner> <title> <base> <head>")
								return
							}
							CreatePR(args[0], args[1], args[2], args[3], args[4])
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
						Name: "create",
						Help: "set the current working issue <repo> <owner> <issuename>",
						Func: func(args []string) {
							if len(args) == 0 || len(args) < 3 {
								fmt.Println("Requires <repo> <owner> <issuename>")
								return
							}
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							if err := CreateIssue(args[0], args[1], args[2]); err != nil {
								color.Red(err.Error())
							} else {
								color.Green("Okay")
							}
						},
					},
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
					command.Command{
						Name: "show",
						Help: "show the current working issue",
						Func: func(args []string) {
							ShowIssue()
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

//CreateIssue ...
func CreateIssue(repo string, owner string, title string) error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}
	}
	githubClient.Issues.List(ctx, true, &github.IssueListOptions{})

	request := &github.IssueRequest{
		Title: &title,
	}
	issue, resp, err := githubClient.Issues.Create(ctx, owner, repo, request)
	if err != nil {
		return err
	}
	fmt.Printf("Github says %d\n", resp.StatusCode)

	localStorage.Github.IssueURL = issue.GetURL()
	storage.Save(localStorage)
	return nil
}

//ShowIssue displays current working issue
func ShowIssue() error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}
	}
	if localStorage.Github.IssueURL != "" {
		fmt.Printf("Working issue at %s\n", localStorage.Github.IssueURL)
	} else {
		color.Red("No working issue set")
	}
	return nil
}

//UnsetIssue from storage
func UnsetIssue() error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}

	}
	localStorage.Github.IssueURL = ""
	return storage.Save(localStorage)
}

//SetIssue in storage
func SetIssue(issueurl string) error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}
	}
	localStorage.Github.IssueURL = issueurl
	return storage.Save(localStorage)
}

//CreatePR makes a new pull request
func CreatePR(owner string, repo string, title string, base string, head string) error {

	pull := github.NewPullRequest{
		Base:  &base,
		Head:  &head,
		Title: &title,
	}

	_, resp, err := githubClient.PullRequests.Create(ctx, owner, repo, &pull)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	fmt.Printf("Github says %d\n", resp.StatusCode)

	return nil
}

//AttachIssuetoPr ...
func AttachIssuetoPr(reponame string, owner string, number string) error {

	if localStorage == nil {
		localStorage, _ = storage.Load()
	}

	if localStorage.Github.IssueURL == "" {
		color.Red("No working issue set...")
		return nil
	}

	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(err)
		return err
	}

	pr, res, err := githubClient.PullRequests.Get(ctx, owner, reponame, num)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Github says %d\n", res.StatusCode)

	appended := fmt.Sprintf("%s\n%s\n", string(pr.GetBody()), localStorage.Github.IssueURL)

	pr, res, err = githubClient.PullRequests.Edit(ctx, owner, reponame, num, &github.PullRequest{Body: &appended})
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("Okay")
	return nil
}
