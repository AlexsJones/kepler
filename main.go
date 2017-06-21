package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlexsJones/kepler/commands"
	"github.com/AlexsJones/kepler/util"
	"github.com/abiosoft/ishell"
	"github.com/dimiro1/banner"
	"github.com/fatih/color"
	"gopkg.in/src-d/go-git.v4"
)

const b string = `
{{ .AnsiColor.Green }} _  _  ____  ____  __    ____  ____
{{ .AnsiColor.Green }}( )/ )( ___)(  _ \(  )  ( ___)(  _ \
{{ .AnsiColor.Green }} )  (  )__)  )___/ )(__  )__)  )   /
{{ .AnsiColor.Green }}(_)\_)(____)(__)  (____)(____)(_)\_)
{{ .AnsiColor.Default }}
{{ .AnsiColor.Default }} Kepler is a simple program for managing submodules
{{ .AnsiColor.Default }} Type 'help' for commands!
{{ .AnsiColor.Default }} Normal shell commands can be used here too e.g. pwd
{{ .AnsiColor.Default }}
`

func main() {
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(b))

	shell := ishell.New()
	shell.SetHomeHistoryPath(".ishell_history")

	shell.AddCmd(&ishell.Cmd{
		Name: "clean",
		Help: "Clean a file in submodule of all references <file> <regex pattern>",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				fmt.Printf("Please provide <file> <regex>\n\n e.g. clean package.json '(?m)[\\r\\n]+^.*substring.*$'\n\n")
				return
			}
			fmt.Printf("Searching %s for %s\n", c.Args[0], c.Args[1])

			valid := regexp.MustCompile(c.Args[1])
			var dryRun bool = true
			reader := bufio.NewReader(os.Stdin)
			color.Red("\nDry run first?[Y/N]\n")
			text, _ := reader.ReadString('\n')
			if strings.Contains(text, "N") {
				dryRun = false
			}
			commands.LoopSubmodules(func(sub *git.Submodule) {
				if err := commands.MatchInfile(sub.Config().Path, c.Args[0], valid, dryRun); err != nil {
					color.Red(err.Error())
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
		util.ShellCommand(strings.Join(arg1.Args, " "), "")
	})
	shell.Run()
}
