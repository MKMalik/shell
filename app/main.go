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

		// check if redirect >> append to file
		redirectAppendStdoutTo := strings.Split(cmd, "1>>")
		redirectAppendStderrTo := strings.Split(cmd, "2>>")

		redirectingAppendStdout := len(redirectAppendStdoutTo) > 1
		redirectingAppendStderr := len(redirectAppendStderrTo) > 1

		if !redirectingAppendStdout {
			redirectAppendStdoutTo = strings.Split(cmd, ">>")
			redirectingAppendStdout = len(redirectAppendStdoutTo) > 1
		}

		// stderr redirect
		if redirectingAppendStderr {
			redirectAppendStdErrToFile(redirectAppendStderrTo)
			continue
		}

		// stdout redirect
		if redirectingAppendStdout {
			redirectAppendStdoutToFile(redirectAppendStdoutTo)
			continue
		}

		// check if redirect > write to file
		redirectWriteStdoutTo := strings.Split(cmd, "1>")
		redirectWriteStderrTo := strings.Split(cmd, "2>")

		redirectingWriteStdout := len(redirectWriteStdoutTo) > 1
		redirectingWriteStderr := len(redirectWriteStderrTo) > 1

		if !redirectingWriteStdout {
			redirectWriteStdoutTo = strings.Split(cmd, ">")
			redirectingWriteStdout = len(redirectWriteStdoutTo) > 1
		}

		// stderr redirect
		if redirectingWriteStderr {
			redirectWriteStdErrToFile(redirectWriteStderrTo)
			continue
		}

		// stdout redirect
		if redirectingWriteStdout {
			redirectWriteStdoutToFile(redirectWriteStdoutTo)
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

func redirectToFile(stdout, stderr, file string, append bool, isStdErr bool) {
	f, err := openFile(file, append)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if isStdErr {
		_, _ = f.WriteString(stderr)
		return
	}

	_, _ = f.WriteString(stdout + "\n")
}

func openFile(name string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return os.Create(name)
}

func redirectAppendStdErrToFile(redirectStderrTo []string) {
	redirectCmd := strings.TrimSpace(redirectStderrTo[0])
	redirectFile := strings.TrimSpace(redirectStderrTo[1])

	stdout, stderr := processCmd(redirectCmd)
	// fmt.Println("Debug: " + stdout + stderr)

	if stdout != "" {
		fmt.Println(stdout)
	}

	redirectToFile(stdout, stderr, redirectFile, true, true)
}

func redirectAppendStdoutToFile(redirectStdoutTo []string) {
	redirectCmd := strings.TrimSpace(redirectStdoutTo[0])
	redirectFile := strings.TrimSpace(redirectStdoutTo[1])

	stdout, stderr := processCmd(redirectCmd)

	// handleRedirectAppendToFile(stdout, redirectFile)
	redirectToFile(stdout, stderr, redirectFile, true, false)

	if stderr != "" {
		fmt.Println(stderr)
	}
}

func redirectWriteStdErrToFile(redirectStderrTo []string) {
	redirectCmd := strings.TrimSpace(redirectStderrTo[0])
	redirectFile := strings.TrimSpace(redirectStderrTo[1])

	stdout, stderr := processCmd(redirectCmd)
	// fmt.Println("Debug: " + stdout + stderr)

	if stdout != "" {
		fmt.Println(stdout)
	}

	// handleRedirectWriteToFile(stderr, redirectFile)
	redirectToFile(stdout, stderr, redirectFile, false, true)
}

func redirectWriteStdoutToFile(redirectStdoutTo []string) {
	redirectCmd := strings.TrimSpace(redirectStdoutTo[0])
	redirectFile := strings.TrimSpace(redirectStdoutTo[1])

	stdout, stderr := processCmd(redirectCmd)

	if stderr != "" {
		fmt.Println(stderr)
	}
	// handleRedirectWriteToFile(stdout, redirectFile)
	redirectToFile(stdout, stderr, redirectFile, false, false)
}

// func handleRedirectWriteToFile(output, redirect string) {
// 	file, err := os.Create(redirect)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	defer file.Close()
// 	_, _ = file.WriteString(output)
// }

// func handleRedirectAppendToFile(output, redirect string) {
// 	file, err := os.OpenFile(redirect, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	defer file.Close()
//
// 	_, _ = file.WriteString(output)
// }

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
