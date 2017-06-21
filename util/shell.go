package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

//ShellCommand ...
func ShellCommand(command string, path string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = path
	stdout, _ := cmd.StdoutPipe()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return
	}
}
