package commands

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"github.com/google/shlex"
)

func ProcessCmd(command string) (string, string) {
	background := false
	if strings.HasSuffix(command, " &") {
		background = true
		command = strings.Split(command, " &")[0]
	}
	fields, _ := shlex.Split(command)

	if len(fields) == 0 {
		return "", ""
	}

	args := fields[1:]
	cmd := handlers.Builtin(fields[0])

	if fn, ok := handlers.Builtins[cmd]; ok {
		return fn(command), ""
	}

	return HandleExternal(fields[0], args, background)
}

func HandleExternal(command string, args []string, background bool) (string, string) {
	found, _ := utils.ScanPath(os.Getenv("PATH"), command)

	if found != nil {
		stdout, stderr, _ := RunCommand(*found, args, background)
		return stdout, stderr
	}

	return "", command + ": command not found\n"
}
