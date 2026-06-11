package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(cmd string, args []string, background bool) (string, string, error) {
	run := exec.Command(cmd, args...)

	if background {
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr

		err := run.Start()
		if err != nil {
			return "", "", err
		}

		pid := run.Process.Pid

		go func() {
			err := run.Wait()

			if err == nil {
				// fmt.Printf("\n[%d]+ Done %s\n", 1, cmd)
			}
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
