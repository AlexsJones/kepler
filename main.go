package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/AlexsJones/kepler/submodules"
	"github.com/abiosoft/ishell"
	"github.com/dimiro1/banner"
)

const b string = `
{{ .AnsiColor.Green }} _  _  ____  ____  __    ____  ____
{{ .AnsiColor.Green }}( )/ )( ___)(  _ \(  )  ( ___)(  _ \
{{ .AnsiColor.Green }} )  (  )__)  )___/ )(__  )__)  )   /
{{ .AnsiColor.Green }}(_)\_)(____)(__)  (____)(____)(_)\_)
{{ .AnsiColor.Default }}
`

func main() {
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(b))

	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "update",
		Help: "Update submodules in directory <path>",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				fmt.Println("Please provide a full path")
				return
			}
			submodules.UpdateSubmodules(c.Args[0])
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "exec",
		Help: "Exec command in submodules <path> <cmd> e.g. exec . \"git reset --hard HEAD\"",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				fmt.Println("Please provide a full path")
				return
			}
			submodules.CommandSubmodules(c.Args[0], c.Args[1])
		},
	})
	shell.Run()
}
