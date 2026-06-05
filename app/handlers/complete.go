package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
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

	if len(fields) >= 3 && fields[1] == "-r" {
		command := fields[2]
		RemoveCompletion(command)
	}

	return ""
}

func GetComplete(cmd string) *string {
	if path, ok := register[cmd]; ok {
		return &path
	}
	return nil
}

var completerFirst = true

type CompletionResult int

const (
	NoCompletion CompletionResult = iota
	Completed
	Handled
)

func RunCompleter(line string) (string, CompletionResult) {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return "", NoCompletion
	}

	arg0 := fields[0]

	path := GetComplete(arg0)
	if path == nil {
		return "", NoCompletion
	}

	arg1 := ""
	arg2 := ""

	if len(fields) >= 2 {
		arg1 = fields[len(fields)-2]
	}

	arg2 = fields[len(fields)-1]

	args := []string{
		fields[0], // command
		arg2,      // current word
		arg1,      // previous word
	}

	cmd := exec.Command(*path, args...)

	cmd.Env = append(
		os.Environ(),
		"COMP_LINE="+line,
		fmt.Sprintf("COMP_POINT=%d", len(line)),
	)

	out, err := cmd.Output()
	if err != nil {
		return "", NoCompletion
	}

	candidates := strings.Fields(string(out))

	if len(candidates) == 0 {
		return "", NoCompletion
	}

	if len(candidates) > 1 {
		current := fields[len(fields)-1]
		lcp := utils.LongestCommonPrefix(candidates)

		if len(lcp) > len(current) {
			fields[len(fields)-1] = lcp

			newLine := strings.Join(fields, " ")

			os.Stdout.WriteString("\r\033[2K$ ")
			os.Stdout.WriteString(newLine)

			completerFirst = true

			return newLine, Completed
		}

		if completerFirst {
			completerFirst = false
			os.Stdout.WriteString("\a")
			return line, Handled
		}

		completerFirst = true

		os.Stdout.WriteString("\r\n")
		os.Stdout.WriteString(strings.Join(candidates, "  "))
		os.Stdout.WriteString("\r\n")
		os.Stdout.WriteString("$ ")
		os.Stdout.WriteString(line)

		return line, Handled
	}
	completerFirst = true

	candidate := candidates[0]

	var newLine string

	if strings.HasSuffix(line, " ") {
		newLine = line + candidate + " "
	} else {
		fields[len(fields)-1] = candidate
		newLine = strings.Join(fields, " ") + " "
	}

	os.Stdout.WriteString("\r\033[2K$ ")
	os.Stdout.WriteString(newLine)

	return newLine, Completed
}

func RemoveCompletion(cmd string) {
	delete(register, cmd)
}
