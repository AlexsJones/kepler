package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/AlexsJones/kepler/commands"
	"github.com/abiosoft/ishell"
	"github.com/dimiro1/banner"
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
	shell.SetHomeHistoryPath(".ishell_history")

	shell.AddCmd(&ishell.Cmd{
		Name: "fix",
		Help: "Fix selected packages to use local links e.g. fix mycompany@git",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Please give a target package string to try to convert to a file link")
				return
			}
			commands.LoopSubmodules(func(sub *git.Submodule) {
				if err := commands.FixLinks(sub.Config().Path, "package.json", c.Args[0]); err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("- Link fixed: %s\n", sub.Config().Path)
				}
			})
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "exec",
		Help: "Exec command in submodules <cmd> e.g. exec \"git reset --hard HEAD\"",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Please provide a command")
				return
			}
			commands.CommandSubmodules(strings.Join(c.Args, " "))
		},
	})
	shell.NotFound(func(arg1 *ishell.Context) {
		// Pass through to bash
		commands.ShellCommand(strings.Join(arg1.Args, " "), "")
	})
	shell.Start()

	shell.Wait()
}
