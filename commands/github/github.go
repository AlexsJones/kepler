//Package github provides a modular way of interacting with github
//This is primary gateway to create/deleting and reviewing both pull requests and issues
package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
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
				Name: "team",
				Help: "team command palette",
				Func: func(args []string) {
					fmt.Println("See help for working with teams")
				},
				SubCommands: []command.Command{
					command.Command{
						Name: "list",
						Help: "List team membership",
						Func: func(args []string) {
							if GithubClient == nil {
								fmt.Println("Please login first...")
								return
							}

							if err := listTeams(); err != nil {
								color.Red(err.Error())
								return
							}
							color.Green("Okay")
						},
					},
					command.Command{
						Name: "set",
						Help: "Set the current team to work with",
						Func: func(args []string) {
							if GithubClient == nil {
								fmt.Println("Please login first...")
								return
							}
							if len(args) == 0 || len(args) < 1 {
								fmt.Println("set the current team id to use <teamid>")
								return
							}
							if err := setTeam(args[0]); err != nil {
								color.Red(err.Error())
							}

							color.Green("Okay")
						},
					},
					command.Command{
						Name: "fetch",
						Help: "Fetch remote team repos",
						Func: func(args []string) {
							if GithubClient == nil {
								fmt.Println("Please login first...")
								return
							}
							if err := fetchTeamRepos(); err != nil {
								color.Red(err.Error())
								return
							}
							color.Green("Okay")
						},
					},
				},
			},

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
							attachIssuetoPr(args[0], args[1], args[2])
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

							if err := createPR(args[0], args[1], args[2], args[3], strings.Join(conc, " ")); err != nil {
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

							if err := createIssue(args[0], args[1], strings.Join(conc, " ")); err != nil {
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
							if err := setIssue(i); err != nil {
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
							if err := unsetIssue(); err != nil {
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
							if err := showIssue(); err != nil {
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
									if err := addRespositoryToPalette(args[0]); err != nil {
										color.Red(err.Error())
										return
									}
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
									if err := deleteRepositoryFromPalette(args[0]); err != nil {
										color.Red(err.Error())
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
									if err := showCurrentIssuePalette(); err != nil {
										color.Red(err.Error())
										return
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
									deleteIssuePalette()
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

					Login()

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
					if err := fetchRepos(); err != nil {
						color.Red(err.Error())
						return
					}
					color.Green("Okay")
				},
			},
		},
	})
}
