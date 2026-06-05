package commands

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RunCommand(found string, args []string, background bool) (string, string, error) {
	run := exec.Command(found, args...)

	if background {
		err := run.Start()
		if err != nil {
			return "", "", err
		}

		pid := run.Process.Pid

		go func() {
			_ = run.Wait()
		}()

		return fmt.Sprintf("[%d] %d\n", 1, pid), "", nil
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	run.Stdout = &stdout
	run.Stderr = &stderr

	err := run.Run()

	return stdout.String(), stderr.String(), err
}
