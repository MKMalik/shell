package commands

import (
	"bytes"
	"os/exec"
)

func RunCommand(found string, args []string) (string, string, error) {
	run := exec.Command(found, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	run.Stdout = &stdout
	run.Stderr = &stderr

	err := run.Run()

	return stdout.String(), stderr.String(), err
}
