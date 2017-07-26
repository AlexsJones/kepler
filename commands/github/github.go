//Package github provides a modular way of interacting with github
//This is primary gateway to create/deleting and reviewing both pull requests and issues
package github

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

//AddCommands for the github module
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
						Help: "attach the current issue to a pr <owner> <reponame> <prnumber>",
						Func: func(args []string) {
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							if len(args) == 0 || len(args) < 3 {
								fmt.Println("set the current working issue in the pr <owner> <reponame> <prnumber>")
								return
							}
							AttachIssuetoPr(args[0], args[1], args[2])
						},
					},
					command.Command{
						Name: "create",
						Help: "create a pr <owner> <repo> <base> <head> <title>",
						Func: func(args []string) {
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							if len(args) == 0 || len(args) < 5 {
								fmt.Println("create a pr <owner> <repo> <base> <head> <title> ")
								return
							}

							var conc []string
							for _, str := range args[4:] {
								conc = append(conc, str)
							}

							if err := CreatePR(args[0], args[1], args[2], args[3], strings.Join(conc, " ")); err != nil {
								color.Red(err.Error())
								return
							}
							color.Green("Okay")
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
						Help: "set the current working issue <owner> <repo> <issuename>",
						Func: func(args []string) {
							if len(args) == 0 || len(args) < 3 {
								fmt.Println("Requires <owner> <repo> <issuename>")
								return
							}
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}

							var conc []string
							for _, str := range args[2:] {
								conc = append(conc, str)
							}

							if err := CreateIssue(args[0], args[1], strings.Join(conc, " ")); err != nil {
								color.Red(err.Error())
							} else {
								color.Green("Okay")
							}
						},
					},
					command.Command{
						Name: "set",
						Help: "set the current working issue <issue number>",
						Func: func(args []string) {
							if len(args) == 0 || len(args) < 1 {
								fmt.Println("Requires <issue number>")
								return
							}
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							i, error := strconv.Atoi(args[0])
							if error != nil {
								color.Red(error.Error())
								return
							}
							if err := SetIssue(i); err != nil {
								color.Red(err.Error())
								return
							}
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
							if err := UnsetIssue(); err != nil {
								color.Red(err.Error())
								return
							}
							color.Green("Okay")
						},
					},
					command.Command{
						Name: "show",
						Help: "show the current working issue",
						Func: func(args []string) {
							if githubClient == nil || localStorage == nil {
								fmt.Println("Please login first...")
								return
							}
							if err := ShowIssue(); err != nil {
								color.Red(err.Error())
								return
							}
							color.Green("Okay")
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
							color.Red(err.Error())
							return
						}
					} else {
						fmt.Print("Access token: ")
						reader := bufio.NewReader(os.Stdin)
						token, _ := reader.ReadString('\n')
						log.Println("Creating new storage object...")
						localStorage = storage.NewStorage()
						localStorage.Github.AccessToken = strings.TrimSpace(token)
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

//CreateIssue creates an issue based on the selected repository
//This will return on success an issue object that is stored in Kepler
func CreateIssue(owner string, repo string, title string) error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}
	}
	fmt.Printf("Owner: %s\n", owner)
	fmt.Printf("Repo: %s\n", repo)
	fmt.Printf("Title: %s\n", title)
	githubClient.Issues.List(ctx, true, &github.IssueListOptions{})

	request := &github.IssueRequest{
		Title: &title,
	}
	issue, resp, err := githubClient.Issues.Create(ctx, owner, repo, request)
	if err != nil {
		return err
	}
	fmt.Printf("Github says %d\n", resp.StatusCode)
	fmt.Printf("%s\n", issue.GetHTMLURL())
	fmt.Printf("Issue status is %s\n", issue.GetState())

	var stIssue storage.Issue
	stIssue.IssueURL = issue.GetHTMLURL()
	stIssue.Owner = owner
	stIssue.Repo = repo
	stIssue.Number = issue.GetNumber()

	localStorage.Github.Issue = append(localStorage.Github.Issue, stIssue)
	storage.Save(localStorage)
	return nil
}

//ShowIssue shows stored issues and highlights the current working issue if set
func ShowIssue() error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}
	}

	for count, currentIssue := range localStorage.Github.Issue {

		issue, _, err := githubClient.Issues.Get(ctx, currentIssue.Owner, currentIssue.Repo, currentIssue.Number)

		if err != nil {
			color.Red(err.Error())
			return err
		}
		if localStorage.Github.CurrentIssue != nil {
			if localStorage.Github.CurrentIssue.IssueURL == currentIssue.IssueURL {
				fmt.Printf("Current issue >>>> ")
			}
		}
		fmt.Printf("%d: issue at %s with status %s\n", count, currentIssue.IssueURL, issue.GetState())

		if len(currentIssue.PullRequests) > 0 {
			fmt.Printf("\n")
			for _, pr := range currentIssue.PullRequests {

				p, _, err := githubClient.PullRequests.Get(ctx, pr.Owner, pr.Repo, pr.Number)
				if err != nil {
					color.Red(err.Error())
					return err
				}
				fmt.Printf("[STATUS:%s]%s/%s  %s base: %s head %s %s\n", p.GetState(), pr.Owner, pr.Repo, p.GetHTMLURL(), pr.Base, pr.Head, pr.Title)

			}
		}
	}
	return nil
}

//UnsetIssue the working issue from storage if set
func UnsetIssue() error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}

	}
	localStorage.Github.CurrentIssue = nil
	return storage.Save(localStorage)
}

//SetIssue in storage using the issue index number
func SetIssue(issueNumber int) error {
	var err error
	if localStorage == nil {
		localStorage, err = storage.Load()
		if err != nil {
			return err
		}
	}

	if issueNumber > len(localStorage.Github.Issue) {
		return errors.New("Out of bounds")
	}

	is := localStorage.Github.Issue[issueNumber]
	if &is == nil {
		return errors.New("No issue pointer")
	}
	localStorage.Github.CurrentIssue = &is
	return storage.Save(localStorage)
}

//CreatePR makes a new pull request with the given criteria
//It returns an error object with nil on success
func CreatePR(owner string, repo string, base string, head string, title string) error {

	fmt.Printf("Owner: %s\n", owner)
	fmt.Printf("Repo: %s\n", repo)
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Base: %s\n", base)
	fmt.Printf("Head: %s\n", head)
	var prbody string
	if localStorage.Github.CurrentIssue.IssueURL != "" {
		fmt.Printf("Attach to the current working issue? (Issue: %s) [Y/N]\n", localStorage.Github.CurrentIssue.IssueURL)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.Contains(response, "Y") {
			prbody = localStorage.Github.CurrentIssue.IssueURL
			fmt.Printf("Body: %s\n", localStorage.Github.CurrentIssue.IssueURL)
		}
	}
	pull := github.NewPullRequest{
		Base:  &base,
		Head:  &head,
		Title: &title,
		Body:  &prbody,
	}
	p, resp, err := githubClient.PullRequests.Create(ctx, owner, repo, &pull)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	fmt.Printf("Github says %d\n", resp.StatusCode)
	fmt.Printf("%s\n", p.GetHTMLURL())
	fmt.Printf("PR status is %s\n", p.GetState())
	storedPr := storage.PullRequest{
		Owner:  owner,
		Repo:   repo,
		Base:   base,
		Head:   head,
		Title:  title,
		Number: p.GetNumber(),
	}
	localStorage.Github.CurrentIssue.PullRequests = append(localStorage.Github.CurrentIssue.PullRequests, storedPr)
	storage.Save(localStorage)
	return nil
}

//AttachIssuetoPr will use the current working issue to attach a new pull request too
func AttachIssuetoPr(owner string, reponame string, number string) error {

	if localStorage == nil {
		localStorage, _ = storage.Load()
	}
	fmt.Printf("Owner: %s\n", owner)
	fmt.Printf("Repo: %s\n", reponame)
	fmt.Printf("Title: %s\n", number)

	if localStorage.Github.CurrentIssue.IssueURL == "" {
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

	appended := fmt.Sprintf("%s\n%s\n", string(pr.GetBody()), localStorage.Github.CurrentIssue.IssueURL)

	pr, res, err = githubClient.PullRequests.Edit(ctx, owner, reponame, num, &github.PullRequest{Body: &appended})
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("Okay")
	return nil
}
