package commands

import "github.com/codecrafters-io/shell-starter-go/app/utils"

func ExecuteRedirect(r utils.Redirect) {
	stdout, stderr := ProcessCmd(r.Cmd, false)

	if r.FD == 1 {
		_ = utils.RedirectToFile(stdout, r.File, r.Append)

		if stderr != "" {
			utils.WriteOutput(stderr)
		}

		return
	}

	_ = utils.RedirectToFile(stderr, r.File, r.Append)

	if stdout != "" {
		utils.WriteOutput(stdout)
	}
}

