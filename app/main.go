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
		// check if cmd contains redirect stdout then redirect without printing
		redirectTo := strings.Split(cmd, "1>")
		if len(redirectTo) < 2 {
			redirectTo = strings.Split(cmd, ">")
		}
		if len(redirectTo) > 1 {
			redirectCmd := redirectTo[0]
			redirect := redirectTo[1]
			output := processCmd(redirectCmd)
			handleRedirectToFile(strings.TrimSpace(output), strings.TrimSpace(redirect))
			continue
		}
		output := processCmd(cmd)
		fmt.Println(output)
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

func processCmd(command string) string {
	fields, _ := shlex.Split(command)
	if len(fields) == 0 {
		return ""
	}
	args := fields[1:]
	switch fields[0] {
	case string(handlers.Echo):
		return handlers.Builtins[handlers.Echo](command)
	case string(handlers.Type):
		return handlers.Builtins[handlers.Type](command)
	case string(handlers.Pwd):
		return handlers.Builtins[handlers.Pwd](command)
	case string(handlers.Cd):
		return handlers.Builtins[handlers.Cd](args[0])
	default:
		// check if exists in PATH as executable
		found, _ := utils.ScanPath(os.Getenv("PATH"), fields[0])
		// if exists and executable then execute passing args if any
		if found != nil {
			output, _ := runCommand(*found, args)
			return output
		}
		// if not: print command not found
		return command + ": command not found"
	}
}

func runCommand(found string, args []string) (string, error) {
	run := exec.Command(found, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	run.Stdout = &stdout
	run.Stderr = &stderr

	err := run.Run()

	output := strings.TrimRight(
		stdout.String()+stderr.String(),
		"\n",
	)

	return output, err
}
