package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
)

func RunCommand(cmd *exec.Cmd, command string, background bool) (string, string) {
	var stdout, stderr bytes.Buffer

	if background {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return "", err.Error()
		}

		jobID := handlers.CurrentJobId
		handlers.CurrentJobId++

		handlers.JobList = append(handlers.JobList, handlers.Job{
			ID:            jobID,
			ProcessID:     cmd.Process.Pid,
			CommandString: command,
			Status:        "Running",
		})

		go func(jobID int) {
			_ = cmd.Wait()

			for i := range handlers.JobList {
				if handlers.JobList[i].ID == jobID {
					handlers.JobList[i].Status = "Done"
					break
				}
			}
		}(jobID)

		return fmt.Sprintf("[%d] %d\n", jobID, cmd.Process.Pid), ""
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// if stderr.Len() > 0 {
		// 	return "", stderr.String()
		// }
		// return "", err.Error()
		return stdout.String(), stderr.String()
	}

	return stdout.String(), stderr.String()
}
