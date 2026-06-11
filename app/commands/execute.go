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

}
