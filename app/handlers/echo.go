package handlers

import (
	"strings"

	"github.com/google/shlex"
)

func HandleEcho(command string) string {
	fields, _ := shlex.Split(command)
	if len(fields) == 0 {
		return ""
	}
	args := fields[1:]
	if len(args) == 0 {
		return ""
	}
	return strings.Join(args, " ") + "\n"
}
