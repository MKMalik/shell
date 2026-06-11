package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"github.com/google/shlex"
)

// func ProcessCmd(command string, background bool) (string, string) {
//
// 	fields, _ := shlex.Split(command)
//
// 	if len(fields) == 0 {
// 		return "", ""
// 	}
//
// 	args := fields[1:]
// 	cmd := handlers.Builtin(fields[0])
//
// 	if fn, ok := handlers.Builtins[cmd]; ok {
// 		return fn(command), ""
// 	}
//
// 	return HandleExternal(fields[0], args, background)
// }

func ProcessCmd(command string, background bool) (string, string) {
	command = strings.TrimSpace(command)

	if command == "" {
		return "", ""
	}

	// Shell expressions
	if containsShellSyntax(command) {
		return HandleShell(command, background)
	}

	fields, err := shlex.Split(command)
	if err != nil || len(fields) == 0 {
		return "", err.Error()
	}

	args := fields[1:]
	cmd := handlers.Builtin(fields[0])

	if fn, ok := handlers.Builtins[cmd]; ok {
		return fn(command), ""
	}

	return HandleExternal(fields[0], args, background)
}

func containsShellSyntax(cmd string) bool {
	ops := []string{
		"&&",
		"||",
		"|",
		";",
		">",
		">>",
		"<",
		"<<",
		"$(",
		"`",
	}

	for _, op := range ops {
		if strings.Contains(cmd, op) {
			return true
		}
	}

	return false
}

func HandleShell(command string, background bool) (string, string) {
	cmd := exec.Command("sh", "-c", command)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if background {
		if err := cmd.Start(); err != nil {
			return "", err.Error()
		}

		go func() {
			cmd.Wait()
			fmt.Print("\r$ ")
		}()

		return "", ""
	}
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", stderr.String()
		}
		return "", err.Error()
	}

	return stdout.String(), stderr.String()
}

func HandleExternal(command string, args []string, background bool) (string, string) {
	found, _ := utils.ScanPath(os.Getenv("PATH"), command)

	if found != nil {
		stdout, stderr, _ := RunCommand(*found, args, background)
		return stdout, stderr
	}

	return "", command + ": command not found\n"
}
