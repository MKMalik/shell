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

	if len(fields) == 0 {
		return "", false
	}

	arg0 := fields[0]

	path := GetComplete(arg0)
	if path == nil {
		return "", false
	}

	arg1 := ""
	arg2 := ""

	if len(fields) >= 2 {
		arg1 = fields[len(fields)-2]
	}

	arg2 = fields[len(fields)-1]

	args := []string{
		fields[0], // git
		arg2,      // get
		arg1,      // remote
	}

	cmd := exec.Command(*path, args...)

	cmd.Env = append(
		os.Environ(),
		"COMP_LINE="+line,
		fmt.Sprintf("COMP_POINT=%d", len(line)),
	)

	out, err := cmd.Output()

	if err != nil {
		return "", false
	}

	candidate := strings.TrimSpace(string(out))
	if candidate == "" {
		return "", false
	}

	var newLine string

	if strings.HasSuffix(line, " ") {
		newLine = line + candidate + " "
	} else {
		fields[len(fields)-1] = candidate
		newLine = strings.Join(fields, " ") + " "
	}

	os.Stdout.WriteString("\r\033[2K$ ")
	os.Stdout.WriteString(newLine)

	return newLine, true
}
