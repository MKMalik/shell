package handlers

import (
	"fmt"
	"strings"

	"github.com/google/shlex"
)

func HandleEcho(command string) {
	fields, _ := shlex.Split(command)
	if len(fields) == 0 {
		return
	}
	args := fields[1:]
	fmt.Println(strings.Join(args, " "))
}
