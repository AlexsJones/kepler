package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/AlexsJones/kepler/commands/node"
	sh "github.com/AlexsJones/kepler/commands/shell"
	"github.com/AlexsJones/kepler/commands/submodules"
	"github.com/AlexsJones/kepler/util"
	"github.com/abiosoft/ishell"
	"github.com/dimiro1/banner"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
)

const b string = `
{{ .AnsiColor.Green }} _  _  ____  ____  __    ____  ____
{{ .AnsiColor.Green }}( )/ )( ___)(  _ \(  )  ( ___)(  _ \
{{ .AnsiColor.Green }} )  (  )__)  )___/ )(__  )__)  )   /
{{ .AnsiColor.Green }}(_)\_)(____)(__)  (____)(____)(_)\_)
{{ .AnsiColor.Default }}
{{ .AnsiColor.Default }} Kepler is a simple program for managing submodules + npm packages
{{ .AnsiColor.Default }} Type 'help' for commands!
{{ .AnsiColor.Default }} Normal shell commands can be used here too e.g. pwd
{{ .AnsiColor.Default }}
`

func main() {
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(b))

	shell := ishell.New()
	shell.SetPrompt("[kepler]>>>")
	shell.SetHomeHistoryPath(".ishell_history")

	shell.AddCmd(&ishell.Cmd{
		Name: "exec",
		Help: "Exec command in submodules <cmd> e.g. exec git reset --hard HEAD",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Please provide a command")
				return
			}
			submodules.CommandSubmodules(strings.Join(c.Args, " "))
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "file",
		Help: "Switch selected packages to use local links e.g. fix mycompany@git",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				fmt.Println("Please give a target package string to try to convert to a file link <prefix> <string> e.g. file ../../ googleremotes.git")
				return
			}
			submodules.LoopSubmodules(func(sub *git.Submodule) {
				if err := node.FixLinks(sub.Config().Path, "package.json", c.Args[0], c.Args[1], false); err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("- Link fixed: %s\n", sub.Config().Path)
				}
			})
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "usage",
		Help: "Find usage in submodules of a certain package e.g. usage mocha",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Find a package usage in submodule package.json e.g. usage mocha")
				return
			}
			submodules.LoopSubmodules(func(sub *git.Submodule) {
				if _, err := node.HasPackage(sub.Config().Path, "package.json", c.Args[0]); err != nil {
				}
			})
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "remove",
		Help: "Remove selected packages that match the <input string> e.g. Google.git",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Please give a target package string to to remove")
				return
			}
			submodules.LoopSubmodules(func(sub *git.Submodule) {
				if err := node.FixLinks(sub.Config().Path, "package.json", "", c.Args[0], true); err != nil {
				} else {
					fmt.Printf("- Removed in: %s\n", sub.Config().Path)
				}
			})
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "purge",
		Help: "Purge all kepler storage",
		Func: func(c *ishell.Context) {

			storage.Delete()
			color.Blue("Deleted local storage")
		},
	})
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
			client := github.NewClient(tc)
			_, _, err = client.Repositories.List(ctx, "", nil)
			if err != nil {
				color.Red("Could not authenticate; please purge and login again")
				color.Red(err.Error())
				return
			}
			color.Green("Authentication Successful.")
		},
	})

	shell.NotFound(func(arg1 *ishell.Context) {
		// Pass through to bash
		sh.ShellCommand(strings.Join(arg1.Args, " "), "", false)
	})

	if len(os.Args) > 1 && os.Args[1] == "unattended" {
		shell.Process(os.Args[2:]...)
	} else {
		shell.Start()
		shell.Wait()
	}
}
