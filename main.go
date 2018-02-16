//Kepler is a tool for helping developers work in the cli with github and other tools
//It's speciality is the management of multiple working issues and threading those together with pull requests
//Ideal audience would be a developer working across multiple repositories
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"github.com/AlexsJones/cli/cli"
	"github.com/AlexsJones/kepler/commands/docker"
	"github.com/AlexsJones/kepler/commands/github"
	"github.com/AlexsJones/kepler/commands/kubebuilder"
	"github.com/AlexsJones/kepler/commands/node"
	"github.com/AlexsJones/kepler/commands/palette"
	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/AlexsJones/kepler/commands/submodules"
	"github.com/AlexsJones/renew"
	"github.com/AlexsJones/renew/fetcher"
	"github.com/dimiro1/banner"
)

var version string = "HEAD"

//Ascii art
const b string = `
{{ .AnsiColor.Green }} _  _  ____  ____  __    ____  ____
{{ .AnsiColor.Green }}( )/ )( ___)(  _ \(  )  ( ___)(  _ \
{{ .AnsiColor.Green }} )  (  )__)  )___/ )(__  )__)  )   /
{{ .AnsiColor.Green }}(_)\_)(____)(__)  (____)(____)(_)\_)
{{ .AnsiColor.Default }}Kepler is a simple program for improving developer workflow.
{{ .AnsiColor.Default }}Type 'help' for commands!
{{ .AnsiColor.Default }}
`

func processVersion() {
	_, filename, _, ok := runtime.Caller(0)
	if ok {

		p := path.Dir(filename)
		if p != "" {
			b, err := storage.Exists(p)
			if err != nil {
				log.Println(err.Error())
			}
			if b {
				cmd := exec.Command("git", "describe", "--always", "--long", "--dirty")
				cmd.Dir = p
				out, err := cmd.Output()
				if err != nil {
					log.Println(err.Error())
				}
				if out != nil {
					version = string(out)
				}
			}
		}
	}

	if version != "" {
		fmt.Printf("Running kepler version %s\n", version)
	}
}
func mainInit() {
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(b))

	processVersion()

	cli := cli.NewCli()

	//Modules to add ----------------------------
	kubebuilder.AddCommands(cli)
	node.AddCommands(cli)
	github.AddCommands(cli)
	submodules.AddCommands(cli)
	storage.AddCommands(cli)
	palette.AddCommands(cli)
	docker.AddCommands(cli)
	//-------------------------------------------
	//Additional commands
	// Only automatically login if there is AccessToken set
	if store := storage.GetInstance(); store.Github.AccessToken != "" {
		github.Login()
	}
	//-------------------------------------------
	cli.Run()
}

func main() {
	stateChange := make(chan renew.StatusCode)

	go func() {
		for {
			select {
			case evt := <-stateChange:
				switch evt {
				case renew.UPDATEFETCHED:
					fmt.Printf("\nA kepler update has been fetched...\n")
					time.Sleep(time.Second * 5)
				case renew.RESTARTING:
					fmt.Println("-----restarting-----")
				}
			}

			time.Sleep(time.Second)
		}
	}()

	fmt.Println("You suck!")
	renew.Run(&renew.Configuration{
		Process:           mainInit,
		StateChange:       stateChange,
		ApplicationGoPath: "github.com/AlexsJones/kepler",
		Fetcher: &fetcher.GithubFetcher{
			Interval:         time.Second * 60,
			GithubRepository: "https://github.com/AlexsJones/kepler.git",
		},
	})
}
