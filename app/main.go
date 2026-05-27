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

		redirectTo := strings.Split(cmd, "1>")

		if len(redirectTo) < 2 {
			redirectTo = strings.Split(cmd, ">")
		}

		if len(redirectTo) > 1 {
			redirectCmd := redirectTo[0]
			redirect := redirectTo[1]

			stdout, stderr := processCmd(redirectCmd)

			handleRedirectToFile(
				strings.TrimSpace(stdout),
				strings.TrimSpace(redirect),
			)

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
		panic(err)
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
