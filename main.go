package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/AlexsJones/kepler/commands/github"
	"github.com/AlexsJones/kepler/commands/node"
	sh "github.com/AlexsJones/kepler/commands/shell"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/AlexsJones/kepler/commands/submodules"
	"github.com/abiosoft/ishell"
	"github.com/dimiro1/banner"
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

	//Modules to add ----------------------------
	node.AddCommands(shell)
	github.AddCommands(shell)
	submodules.AddCommands(shell)
	storage.AddCommands(shell)
	//-------------------------------------------
	shell.NotFound(func(arg1 *ishell.Context) {
		sh.ShellCommand(strings.Join(arg1.Args, " "), "", false)
	})

	if len(os.Args) > 1 && os.Args[1] == "unattended" {
		shell.Process(os.Args[2:]...)
	} else {
		shell.Start()
		shell.Wait()
	}
}
