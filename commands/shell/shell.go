//Package sh acts as a wrapper around system commands
package sh

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

//ShellCommand lets you perform a shell command in bash
//It also captures stderr and stdout into Kepler
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
