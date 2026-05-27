package handlers

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleType(cmd string) string {
	fields := strings.Fields(cmd)
	if len(fields) < 2 {
		return ""
	}
	arg := fields[1]
	if !IsBuiltin(arg) {
		paths := os.Getenv("PATH")
		_, full := utils.ScanPath(paths, arg)
		if full != nil {
			return arg + " is " + *full

		}
		return arg + ": not found"
	}

	return arg + " is a shell builtin"
}
