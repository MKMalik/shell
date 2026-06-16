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

		handlers.JobMutex.Lock()

		jobID := handlers.GetNextJobID()

		handlers.JobList = append(handlers.JobList, handlers.Job{
			ID:            jobID,
			ProcessID:     cmd.Process.Pid,
			CommandString: command,
			Status:        "Running",
		})

		handlers.JobMutex.Unlock()

		go func(id int) {

			_ = cmd.Wait()

			handlers.JobMutex.Lock()
			defer handlers.JobMutex.Unlock()

			for i := range handlers.JobList {
				if handlers.JobList[i].ID == id {
					handlers.JobList[i].Status = "Done"
					break
				}
			}

		}(jobID)

		return fmt.Sprintf("[%d] %d\n", jobID, cmd.Process.Pid), ""
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return stdout.String(), stderr.String()
	}

	return stdout.String(), stderr.String()
}
