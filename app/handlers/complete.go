package handlers

import (
	"fmt"
	"strings"

	"github.com/google/shlex"
)

var register = map[string]string{}

func HandleComplete(cmd string) string {
	fields, _ := shlex.Split(cmd)
	if len(fields) == 0 {
		return ""
	}

	// complete -C <script> <cmd>
	if len(fields) >= 3 && fields[1] == "-C" {
		script := strings.Trim(fields[2], "'")
		command := fields[3]

		register[command] = script
		return ""
	}

	// complete -p <cmd>
	if len(fields) == 3 && fields[1] == "-p" {
		command := fields[2]

		script, ok := register[command]
		if !ok {
			return fmt.Sprintf("complete: %s: no completion specification\n", command)
		}

		return fmt.Sprintf("complete -C '%s' %s\n", script, command)
	}

	return ""
}
