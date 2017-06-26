package sh

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
)

//AddCommands ...
func AddCommands(shell *ishell.Shell) string {
	shell.NotFound(func(arg1 *ishell.Context) {
		ShellCommand(strings.Join(arg1.Args, " "), "", false)
	})
	return "sh"
}

//ShellCommand ...
func ShellCommand(command string, path string, validated bool) {
	cmd := exec.Command("bash", "-c", command)
	if path != "" {
		cmd.Dir = path
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	scanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()
	go func() {
		for errScanner.Scan() {
			fmt.Printf("%s\n", errScanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if validated {
			color.Green("[%s]OK\n", path)
		}
	}

}
