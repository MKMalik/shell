package commands

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func ExecuteCommand(cmd string) {
	if strings.TrimSpace(cmd) == "exit" {
		os.Exit(0)
	}

	redirect := utils.ParseRedirect(cmd)

	if redirect.Valid {
		ExecuteRedirect(redirect)
		return
	}

	stdout, stderr := ProcessCmd(cmd)

	utils.WriteOutput(stdout)
	utils.WriteOutput(stderr)
}
