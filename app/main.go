package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"github.com/google/shlex"
)

func main() {
	for i := range handlers.Builtins {
		handlers.BuiltinNames[i] = struct{}{}
	}
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")

		cmd, err := scanner.ReadString('\n')
		if err != nil {
			panic(err)
		}

		cmd = strings.TrimSpace(cmd)

		if cmd == "exit" {
			break
		}

		redirectStdoutTo := strings.Split(cmd, "1>")
		redirectStderrTo := strings.Split(cmd, "2>")

		redirectingStdout := len(redirectStdoutTo) > 1
		redirectingStderr := len(redirectStderrTo) > 1

		if !redirectingStdout {
			redirectStdoutTo = strings.Split(cmd, ">")
			redirectingStdout = len(redirectStdoutTo) > 1
		}

		// stderr redirect
		if redirectingStderr {
			redirectCmd := strings.TrimSpace(redirectStderrTo[0])
			redirectFile := strings.TrimSpace(redirectStderrTo[1])

			stdout, stderr := processCmd(redirectCmd)
			// fmt.Println("Debug: " + stdout + stderr)

			if stdout != "" {
				fmt.Println(stdout)
			}

			handleRedirectToFile(stderr, redirectFile)

			continue
		}

		// stdout redirect
		if redirectingStdout {
			redirectCmd := strings.TrimSpace(redirectStdoutTo[0])
			redirectFile := strings.TrimSpace(redirectStdoutTo[1])

			stdout, stderr := processCmd(redirectCmd)

			handleRedirectToFile(stdout, redirectFile)

			if stderr != "" {
				fmt.Println(stderr)
			}

			continue
		}

		stdout, stderr := processCmd(cmd)

		if stdout != "" {
			fmt.Println(stdout)
		}

		if stderr != "" {
			fmt.Println(stderr)
		}
	}
}

func handleRedirectToFile(output, redirect string) {
	file, err := os.Create(redirect)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	file.WriteString(output)
}

func processCmd(command string) (string, string) {
	fields, _ := shlex.Split(command)

	if len(fields) == 0 {
		return "", ""
	}

	args := fields[1:]

	switch fields[0] {
	case string(handlers.Echo):
		return handlers.Builtins[handlers.Echo](command), ""

	case string(handlers.Type):
		return handlers.Builtins[handlers.Type](command), ""

	case string(handlers.Pwd):
		return handlers.Builtins[handlers.Pwd](command), ""

	case string(handlers.Cd):
		return handlers.Builtins[handlers.Cd](args[0]), ""

	default:
		found, _ := utils.ScanPath(os.Getenv("PATH"), fields[0])

		if found != nil {
			stdout, stderr, _ := runCommand(*found, args)
			return stdout, stderr
		}

		return "", command + ": command not found"
	}
}
func runCommand(found string, args []string) (string, string, error) {
	run := exec.Command(found, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	run.Stdout = &stdout
	run.Stderr = &stderr

	err := run.Run()

	return strings.TrimRight(stdout.String(), "\n"),
		strings.TrimRight(stderr.String(), "\n"),
		err
}
