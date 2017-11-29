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
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//GithubClient is the global github interface
var GithubClient *github.Client

//Ctx is the github oauth context
var Ctx context.Context

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
							if GithubClient == nil {
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
							if GithubClient == nil {
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
				Help: "Issue commands",
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
							if GithubClient == nil {
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
							if GithubClient == nil {
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
							if GithubClient == nil {
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
							if GithubClient == nil {
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
					command.Command{
						Name: "palette",
						Help: "Manipulate the issue palette of working repos",
						Func: func(args []string) {
							fmt.Println("Please run palette commands from your meta repo working directory")
						},
						SubCommands: []command.Command{
							command.Command{
								Name: "add",
								Help: "Add a repository to the palette as part of current working issue by name <name>",
								Func: func(args []string) {
									if len(args) == 0 || len(args) < 1 {
										fmt.Println("Requires <issue number>")
										return
									}
									if GithubClient == nil {
										fmt.Println("Please login first...")
										return
									}
									if storage.GetInstance().Github.CurrentIssue == nil {
										fmt.Println("There is no working issue set; set with github issue set")
										return
									}
									if _, err := os.Stat(args[0]); os.IsNotExist(err) {
										color.Red(fmt.Sprintf("The named repo %s does not exist as a sub directory of the current working directory", args[0]))
										return
									}
									dir, err := os.Getwd()
									if err != nil {
										log.Fatal(err)
									}
									p := path.Join(dir, args[0])
									storage.GetInstance().Github.CurrentIssue.Palette[args[0]] = p
									storage.GetInstance().Save()
									color.Green("Okay")
								},
							},
							command.Command{
								Name: "remove",
								Help: "Remove a repository from the palette as part of the current working issue by name <name>",
								Func: func(args []string) {
									if len(args) == 0 || len(args) < 1 {
										fmt.Println("Requires <issue number>")
										return
									}
									if GithubClient == nil {
										fmt.Println("Please login first...")
										return
									}
									if storage.GetInstance().Github.CurrentIssue == nil {
										fmt.Println("There is no working issue set; set with github issue set")
										return
									}
									found := false
									for k := range storage.GetInstance().Github.CurrentIssue.Palette {
										if strings.Compare(k, args[0]) == 0 {
											found = true
											delete(storage.GetInstance().Github.CurrentIssue.Palette, k)
											storage.GetInstance().Save()
										}
									}
									if found != true {
										color.Red(fmt.Sprintf("There was no repo matching the name %s in the palette", args[0]))
										return
									}
									color.Green("Okay")
								},
							},
							command.Command{
								Name: "show",
								Help: "Show repositories in the palette as part of the current working issue",
								Func: func(args []string) {

									if GithubClient == nil {
										fmt.Println("Please login first...")
										return
									}
									if storage.GetInstance().Github.CurrentIssue == nil {
										fmt.Println("There is no working issue set; set with github issue set")
										return
									}
									for k, v := range storage.GetInstance().Github.CurrentIssue.Palette {
										cmd := exec.Command("git", "branch")
										cmd.Dir = v
										out, err := cmd.Output()
										if err != nil {
											color.Red(err.Error())
											return
										}
										ar := strings.Split(string(out), " ")
										trimmed := strings.TrimSuffix(string(ar[1]), "\n")
										trimmed = strings.TrimPrefix(trimmed, "*")
										trimmed = strings.TrimSpace(trimmed)
										fmt.Println(fmt.Sprintf("Name: %s Branch: %s Path: %s", k, trimmed, v))
									}
									color.Green("Okay")
								},
							},
							command.Command{
								Name: "delete",
								Help: "Delete all repositories in the palette as part of the current working issue",
								Func: func(args []string) {

									if GithubClient == nil {
										fmt.Println("Please login first...")
										return
									}
									if storage.GetInstance().Github.CurrentIssue == nil {
										fmt.Println("There is no working issue set; set with github issue set")
										return
									}
									storage.GetInstance().Github.CurrentIssue.Palette = make(map[string]string)
									color.Green("Okay")
								},
							},
						},
					},
				},
			},
			command.Command{
				Name: "login",
				Help: "use an access token to login to github",
				Func: func(args []string) {

					fmt.Print("Access token: ")
					reader := bufio.NewReader(os.Stdin)
					token, _ := reader.ReadString('\n')
					log.Println("Creating new storage tooken...")
					storage.GetInstance().Github.AccessToken = strings.TrimSpace(token)
					storage.GetInstance().Save()

					Ctx = context.Background()
					ts := oauth2.StaticTokenSource(
						&oauth2.Token{AccessToken: storage.GetInstance().Github.AccessToken},
					)
					tc := oauth2.NewClient(Ctx, ts)
					GithubClient = github.NewClient(tc)
					_, _, err := GithubClient.Repositories.List(Ctx, "", nil)
					if err != nil {
						color.Red("Could not authenticate; please purge and login again")
						color.Red(err.Error())
						return
					}
					color.Green("Authentication Successful.")
				},
			},
			command.Command{
				Name: "fetch",
				Help: "fetch remote repos",
				Func: func(args []string) {
					if GithubClient == nil {
						fmt.Println("Please login first...")
						return
					}
					if err := FetchRepos(); err != nil {
						color.Red(err.Error())
						return
					}
					color.Green("Okay")
				},
			},
		},
	})
}

//CreateIssue creates an issue based on the selected repository
//This will return on success an issue object that is stored in Kepler
func CreateIssue(owner string, repo string, title string) error {
	var err error
	fmt.Printf("Owner: %s\n", owner)
	fmt.Printf("Repo: %s\n", repo)
	fmt.Printf("Title: %s\n", title)
	GithubClient.Issues.List(Ctx, true, &github.IssueListOptions{})

	request := &github.IssueRequest{
		Title: &title,
	}
	issue, resp, err := GithubClient.Issues.Create(Ctx, owner, repo, request)
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
	stIssue.Palette = make(map[string]string)

	storage.GetInstance().Github.Issue = append(storage.GetInstance().Github.Issue, stIssue)
	storage.GetInstance().Save()
	return nil
}

//ShowIssue shows stored issues and highlights the current working issue if set
func ShowIssue() error {

	if len(storage.GetInstance().Github.Issue) == 0 {
		return errors.New("No issue set")
	}
	for count, currentIssue := range storage.GetInstance().Github.Issue {

		issue, _, err := GithubClient.Issues.Get(Ctx, currentIssue.Owner, currentIssue.Repo, currentIssue.Number)

		if err != nil {
			color.Red(err.Error())
			return err
		}
		if storage.GetInstance().Github.CurrentIssue != nil {
			if storage.GetInstance().Github.CurrentIssue.IssueURL == currentIssue.IssueURL {
				fmt.Printf("Current issue >>>> ")
			}
		}
		fmt.Printf("%d: issue at %s with status %s\n", count, currentIssue.IssueURL, issue.GetState())

		if len(currentIssue.PullRequests) > 0 {
			fmt.Printf("\n")
			for _, pr := range currentIssue.PullRequests {

				p, _, err := GithubClient.PullRequests.Get(Ctx, pr.Owner, pr.Repo, pr.Number)
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

	if storage.GetInstance().Github.CurrentIssue == nil {
		return errors.New("No issue to unset")
	}
	storage.GetInstance().Github.CurrentIssue = nil
	return storage.GetInstance().Save()
}

//SetIssue in storage using the issue index number
func SetIssue(issueNumber int) error {

	if issueNumber > len(storage.GetInstance().Github.Issue) {
		return errors.New("Out of bounds")
	}

	is := storage.GetInstance().Github.Issue[issueNumber]
	if &is == nil {
		return errors.New("No issue pointer")
	}
	storage.GetInstance().Github.CurrentIssue = &is
	return storage.GetInstance().Save()
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
	if storage.GetInstance().Github.CurrentIssue.IssueURL != "" {
		fmt.Printf("Attach to the current working issue? (Issue: %s) [Y/N]\n", storage.GetInstance().Github.CurrentIssue.IssueURL)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.Contains(response, "Y") {
			prbody = storage.GetInstance().Github.CurrentIssue.IssueURL
			fmt.Printf("Body: %s\n", storage.GetInstance().Github.CurrentIssue.IssueURL)
		}
	}
	pull := github.NewPullRequest{
		Base:  &base,
		Head:  &head,
		Title: &title,
		Body:  &prbody,
	}
	p, resp, err := GithubClient.PullRequests.Create(Ctx, owner, repo, &pull)
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
	storage.GetInstance().Github.CurrentIssue.PullRequests = append(storage.GetInstance().Github.CurrentIssue.PullRequests, storedPr)
	storage.GetInstance().Save()
	return nil
}

//AttachIssuetoPr will use the current working issue to attach a new pull request too
func AttachIssuetoPr(owner string, reponame string, number string) error {

	fmt.Printf("Owner: %s\n", owner)
	fmt.Printf("Repo: %s\n", reponame)
	fmt.Printf("Title: %s\n", number)

	if storage.GetInstance().Github.CurrentIssue.IssueURL == "" {
		color.Red("No working issue set...")
		return nil
	}

	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(err)
		return err
	}

	pr, res, err := GithubClient.PullRequests.Get(Ctx, owner, reponame, num)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Github says %d\n", res.StatusCode)

	appended := fmt.Sprintf("%s\n%s\n", string(pr.GetBody()), storage.GetInstance().Github.CurrentIssue.IssueURL)

	pr, res, err = GithubClient.PullRequests.Edit(Ctx, owner, reponame, num, &github.PullRequest{Body: &appended})
	if err != nil {
		fmt.Println(err)
		return err
	}
	color.Green("Okay")
	return nil
}

//FetchRepos into the current working directory
func FetchRepos() error {

	var count = 0
	var repoList = make(map[string]string)

	opts := github.RepositoryListOptions{}

	opts.PerPage = 20
	for {
		opts.Page = count
		repos, resp, err := GithubClient.Repositories.List(Ctx, "", &opts)
		if err != nil {
			return err
		}
		if len(repos) == 0 || err != nil || resp.StatusCode != 200 {
			break
		}
		log.Printf("Fetched page %d from github\n", count)
		count++

		for _, repo := range repos {
			repoList[repo.GetName()] = repo.GetSSHURL()
		}
	}

	for k, v := range repoList {
		fmt.Printf("%s -> %s\n", k, v)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Fetch from remotes?(Y/N): ")
	text, _ := reader.ReadString('\n')
	if strings.Contains(text, "Y") {

		for name, repo := range repoList {
			fmt.Printf("Fetching %s\n", name)
			out, err := exec.Command("git", "clone", fmt.Sprintf("%s", repo)).Output()
			if err != nil {
				color.Red(fmt.Sprintf("%s %s", string(out), err.Error()))
			}
			color.Green(fmt.Sprintf("Fetched %s\n", name))
			time.Sleep(time.Second)
		}
	}
	return nil
}
