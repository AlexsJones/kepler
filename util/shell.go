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
	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		go func() {
			for errScanner.Scan() {
				fmt.Printf("%s\n", errScanner.Text())
			}
		}()
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		go func() {
			for errScanner.Scan() {
				fmt.Printf("%s\n", errScanner.Text())
			}
		}()
		return
	}
}
