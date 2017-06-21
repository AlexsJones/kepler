package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/AlexsJones/kepler/submodules"
	"github.com/AlexsJones/kepler/util"
	"github.com/abiosoft/ishell"
	"github.com/dimiro1/banner"
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
		Name: "exec",
		Help: "Exec command in submodules <cmd> e.g. exec \"git reset --hard HEAD\"",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Please provide a command")
				return
			}
			submodules.CommandSubmodules(strings.Join(c.Args, " "))
		},
	})
	shell.NotFound(func(arg1 *ishell.Context) {

		util.ShellCommand(strings.Join(arg1.Args, " "), "")
	})
	shell.Run()
}
