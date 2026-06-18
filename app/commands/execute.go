package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func ExecuteCommand(cmd string) {
	if strings.TrimSpace(cmd) == "exit" {
		handlers.HandleExit()
	}

	cmd = ExpandVars(cmd, handlers.GetDeclaredMap())

	background := false

	if strings.HasSuffix(cmd, "&") {
		background = true
		cmd = strings.TrimSpace(strings.TrimSuffix(cmd, "&"))
	}

	redirect := utils.ParseRedirect(cmd)

	if redirect.Valid {
		ExecuteRedirect(redirect)
		return
	}

	stdout, stderr := ProcessCmd(cmd, background)

	utils.WriteOutput(stdout)
	utils.WriteOutput(stderr)

	if msg := handlers.ReapJobs(); msg != "" {
		fmt.Print(msg)
	}
}

func ExpandVars(input string, vars map[string]string) string {
	var out strings.Builder

	for i := 0; i < len(input); i++ {
		if input[i] != '$' {
			out.WriteByte(input[i])
			continue
		}

		i++

		if i >= len(input) {
			out.WriteByte('$')
			break
		}

		var name string

		if input[i] == '{' {
			// ${VAR}
			i++
			start := i

			for i < len(input) && input[i] != '}' {
				i++
			}

			name = input[start:i]

		} else {
			// $VAR
			start := i

			for i < len(input) && isVarChar(input[i]) {
				i++
			}

			name = input[start:i]
			i--
		}

		if val, ok := vars[name]; ok {
			out.WriteString(val)
		}
	}

	return out.String()
}

func isVarChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}
