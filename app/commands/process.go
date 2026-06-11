package commands

import (
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

	// shell operators: && || | ; > < ...
	if containsShellSyntax(command) {
		return HandleShell(command, background)
	}

	fields, err := shlex.Split(command)
	if err != nil {
		return "", err.Error()
	}

	if len(fields) == 0 {
		return "", ""
	}

	if fn, ok := handlers.Builtins[handlers.Builtin(fields[0])]; ok {
		return fn(command), ""
	}

	return HandleExternal(
		fields[0],
		fields[1:],
		background,
	)
}

func containsShellSyntax(cmd string) bool {
	ops := []string{
		"&&",
		"||",
		"|",
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
	return RunCommand(cmd, command, background)
}

func HandleExternal(command string, args []string, background bool) (string, string) {
	found, _ := utils.ScanPath(os.Getenv("PATH"), command)

	if found == nil {
		return "", command + ": command not found\n"
	}

	cmd := exec.Command(*found, args...)

	return RunCommand(
		cmd,
		command+" "+strings.Join(args, " "),
		background,
	)
}
