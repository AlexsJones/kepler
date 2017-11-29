package palette

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/cli/command"
	"github.com/AlexsJones/kepler/commands/github"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/fatih/color"
)

//AddCommands for the palette module
func AddCommands(cli *cli.Cli) {

	cli.AddCommand(command.Command{
		Name: "palette",
		Help: "Issue palette that controls repos can be used from here",
		Func: func(args []string) {
			fmt.Println("See help for working with palette")
		},
		SubCommands: []command.Command{
			command.Command{
				Name: "branch",
				Help: "switch branches or create if they don't exist for working issue palette repos <branchname>",
				Func: func(args []string) {
					fmt.Println("See help for working with palette branch")
				},
				SubCommands: []command.Command{
					command.Command{
						Name: "push",
						Help: "For pushing the local branches to new/existing remotes",
						Func: func(args []string) {
							if github.GithubClient == nil {
								fmt.Println("Please login first...")
								return
							}
							if len(args) == 0 || len(args) < 1 {
								fmt.Println("provide the branch name to switch repo in the palette to <branchname>")
								return
							}
							for k, v := range storage.GetInstance().Github.CurrentIssue.Palette {
								if _, err := os.Stat(v); os.IsNotExist(err) {
									color.Red(fmt.Sprintf("Warning the repo %s does not exist at the path %s, removing from the palette\n", k, v))
									delete(storage.GetInstance().Github.CurrentIssue.Palette, k)
									storage.GetInstance().Save()
								} else if _, err := os.Stat(path.Join(v, ".git")); os.IsNotExist(err) {
									color.Red(fmt.Sprintf("%s .git directory does not exist removing from the palette\n", k))
									delete(storage.GetInstance().Github.CurrentIssue.Palette, k)
									storage.GetInstance().Save()
								} else {
									color.Green(fmt.Sprintf("Pushing %s branches to remote %s:%s", k, args[0], args[0]))
									cmd := exec.Command("git", "push", "origin", fmt.Sprintf("%s:%s", args[0], args[0]))
									cmd.Dir = v
									_, err := cmd.Output()
									if err != nil {
										color.Red(err.Error())
										break
									}
								}
							}
						},
					},
					command.Command{
						Name: "local",
						Help: "For switching local branches on palette repos",
						Func: func(args []string) {
							if github.GithubClient == nil {
								fmt.Println("Please login first...")
								return
							}
							if len(args) == 0 || len(args) < 1 {
								fmt.Println("provide the branch name to switch repo in the palette to <branchname>")
								return
							}
							for k, v := range storage.GetInstance().Github.CurrentIssue.Palette {
								if _, err := os.Stat(v); os.IsNotExist(err) {
									color.Red(fmt.Sprintf("Warning the repo %s does not exist at the path %s, removing from the palette\n", k, v))
									delete(storage.GetInstance().Github.CurrentIssue.Palette, k)
									storage.GetInstance().Save()
								} else if _, err := os.Stat(path.Join(v, ".git")); os.IsNotExist(err) {
									color.Red(fmt.Sprintf("%s .git directory does not exist removing from the palette\n", k))
									delete(storage.GetInstance().Github.CurrentIssue.Palette, k)
									storage.GetInstance().Save()
								} else {
									color.Green(fmt.Sprintf("Switching %s to branch %s", k, args[0]))
									cmd := exec.Command("git", "branch", args[0])
									cmd.Dir = v
									_, err := cmd.Output()
									if err != nil {
										color.Red(err.Error())
										break
									}
									cmd = exec.Command("git", "checkout", args[0])
									cmd.Dir = v
									_, err = cmd.Output()
									if err != nil {
										color.Red(err.Error())
										break
									}
								}
							}
						},
					},
				},
			},
			command.Command{
				Name: "show",
				Help: "Show repositories in the palette as part of the current working issue",
				Func: func(args []string) {

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
		},
	})
}
