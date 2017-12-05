package github

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//Login function to github
func Login() {
	color.Blue("Logging into github...")
	if storage.GetInstance().Github.AccessToken == "" {
		fmt.Print("Access token: ")
		reader := bufio.NewReader(os.Stdin)
		token, _ := reader.ReadString('\n')

		storage.GetInstance().Github.AccessToken = strings.TrimSpace(token)
		storage.GetInstance().Save()
	}
	Ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: storage.GetInstance().Github.AccessToken},
	)
	tc := oauth2.NewClient(Ctx, ts)
	GithubClient = github.NewClient(tc)
	_, _, err := GithubClient.Repositories.List(Ctx, "", nil)
	if err != nil {
		color.Red("Could not authenticate with Github; please purge and login again")
		color.Red(err.Error())
		return
	}
	color.Green("Authentication with Github Successful.")
}

//CreateIssue creates an issue based on the selected repository
//This will return on success an issue object that is stored in Kepler
func createIssue(owner string, repo string, title string) error {
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
func showIssue() error {

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
func unsetIssue() error {

	if storage.GetInstance().Github.CurrentIssue == nil {
		return errors.New("No issue to unset")
	}
	storage.GetInstance().Github.CurrentIssue = nil
	return storage.GetInstance().Save()
}

//SetIssue in storage using the issue index number
func setIssue(issueNumber int) error {

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
func createPR(owner string, repo string, base string, head string, title string) error {

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
func attachIssuetoPr(owner string, reponame string, number string) error {

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

func fetchRepoList(repoList map[string]string) error {

	for k, v := range repoList {
		fmt.Printf("%s -> %s\n", k, v)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Fetch from remotes?(Y/N): ")
	text, _ := reader.ReadString('\n')
	if strings.Contains(text, "Y") {
		_, er := os.Stat(".gitmodules")
		isMetaRepo := !os.IsNotExist(er)
		for name, repo := range repoList {
			if _, er := os.Stat(name); !os.IsNotExist(er) {
				color.Blue("Already have %s", name)
				continue
			}
			var out []byte
			var err error
			fmt.Printf("Fetching %s\n", name)
			if isMetaRepo {
				out, err = exec.Command("git", "submodule", "add", fmt.Sprintf("%s", repo)).Output()
			} else {
				out, err = exec.Command("git", "clone", fmt.Sprintf("%s", repo)).Output()
			}
			if err != nil {
				color.Red(fmt.Sprintf("%s %s", string(out), err.Error()))
			} else {
				color.Green(fmt.Sprintf("Fetched %s\n", name))
			}
			time.Sleep(time.Second)
		}
	}
	return nil
}

//FetchTeamRepos ...
func fetchTeamRepos() error {

	var repoList = make(map[string]string)

	repos, _, err := GithubClient.Organizations.ListTeamRepos(Ctx, storage.GetInstance().Github.TeamID, &github.ListOptions{})
	if err != nil {
		return err
	}
	for _, repo := range repos {

		repoList[repo.GetName()] = repo.GetSSHURL()

	}

	return fetchRepoList(repoList)
}

//FetchRepos into the current working directory
func fetchRepos() error {

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
	return fetchRepoList(repoList)
}

func setTeam(team string) error {
	i, err := strconv.Atoi(team)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	storage.GetInstance().Github.TeamID = i
	storage.GetInstance().Save()
	return nil
}
func listTeams() error {
	teams, _, err := GithubClient.Organizations.ListTeams(Ctx, "SeedJobs", &github.ListOptions{})
	if err != nil {
		return err
	}
	currentTeamID := storage.GetInstance().Github.TeamID
	for _, t := range teams {
		if currentTeamID != 0 && currentTeamID == t.GetID() {
			fmt.Printf("Name: %s -- ID: %d [Currently set team]\n", t.GetName(), t.GetID())
		} else {
			fmt.Printf("Name: %s -- ID: %d\n", t.GetName(), t.GetID())
		}
	}
	return nil
}
