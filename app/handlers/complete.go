package handlers

import (
	"fmt"
	"os"
	"os/exec"
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

func GetComplete(cmd string) *string {
	if path, ok := register[cmd]; ok {
		return &path
	}
	return nil
}

func RunCompleter(line string) (string, bool) {
	fields := strings.Fields(line)

	if len(fields) != 1 {
		return "", false
	}

	cmd := fields[0]

	path := GetComplete(cmd)
	if path == nil {
		return "", false
	}

	out, err := exec.Command(*path).Output()
	if err != nil {
		return "", false
	}

	candidate := strings.TrimSpace(string(out))
	if candidate == "" {
		return "", false
	}

	newLine := line + candidate + " "

	os.Stdout.WriteString("\r\033[2K$ ")
	os.Stdout.WriteString(strings.Split(newLine, "\n")[0]) // single line

	return newLine, true
}
