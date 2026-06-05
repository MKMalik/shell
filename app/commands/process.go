package commands

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"github.com/google/shlex"
)

func ProcessCmd(command string) (string, string) {
	fields, _ := shlex.Split(command)

	if len(fields) == 0 {
		return "", ""
	}

	args := fields[1:]
	cmd := handlers.Builtin(fields[0])

	if fn, ok := handlers.Builtins[cmd]; ok {
		return fn(command), ""
	}

	return HandleExternal(fields[0], args)
}

func HandleExternal(command string, args []string) (string, string) {
	found, _ := utils.ScanPath(os.Getenv("PATH"), command)

	if found != nil {
		stdout, stderr, _ := RunCommand(*found, args)
		return stdout, stderr
	}

	return "", command + ": command not found\n"
}
