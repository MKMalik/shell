package commands

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"github.com/google/shlex"
)

// func ProcessCmd(command string, background bool) (string, string) {
//
// 	fields, _ := shlex.Split(command)
//
// 	if len(fields) == 0 {
// 		return "", ""
// 	}
//
// 	args := fields[1:]
// 	cmd := handlers.Builtin(fields[0])
//
// 	if fn, ok := handlers.Builtins[cmd]; ok {
// 		return fn(command), ""
// 	}
//
// 	return HandleExternal(fields[0], args, background)
// }

func ProcessCmd(command string, background bool) (string, string) {
	command = strings.TrimSpace(command)

	if command == "" {
		return "", ""
	}

	if containsPipe(command) {
		return HandlePipeline(command, background)
	}

	// shell operators: && || | ; > < ...
	if containsShellSyntax(command) {
		return HandleShell(command, background)
	}

	fields, err := shlex.Split(command)
	if err != nil {
		return "", err.Error()
	}

	if len(fields) == 0 {
		return "", ""
	}

	if fn, ok := handlers.Builtins[handlers.Builtin(fields[0])]; ok {
		return fn(command), ""
	}

	return HandleExternal(
		fields[0],
		fields[1:],
		background,
	)
}

func containsPipe(cmd string) bool {
	return !strings.Contains(cmd, "||") && strings.Contains(cmd, "|")
}

func containsShellSyntax(cmd string) bool {
	ops := []string{
		"&&",
		"||",
	}

	for _, op := range ops {
		if strings.Contains(cmd, op) {
			return true
		}
	}

	return false
}

func HandleShell(command string, background bool) (string, string) {
	cmd := exec.Command("sh", "-c", command)
	return RunCommand(cmd, command, background)
}

func HandlePipeline(command string, background bool) (string, string) {
	parts := strings.SplitN(command, "|", 2)

	left := strings.Fields(strings.TrimSpace(parts[0]))
	right := strings.Fields(strings.TrimSpace(parts[1]))

	cmd1 := exec.Command(left[0], left[1:]...)
	cmd2 := exec.Command(right[0], right[1:]...)

	reader, writer := io.Pipe()

	cmd1.Stdout = writer
	cmd1.Stderr = os.Stderr

	cmd2.Stdin = reader

	var out bytes.Buffer

	cmd2.Stdout = io.MultiWriter(
		LineWriter{os.Stdout},
		&out,
	)
	cmd2.Stderr = os.Stderr

	if err := cmd1.Start(); err != nil {
		return "", err.Error()
	}

	if err := cmd2.Start(); err != nil {
		return "", err.Error()
	}

	// drain output while cmd2 runs
	cmd2Done := make(chan error, 1)

	go func() {
		cmd2Done <- cmd2.Wait()
	}()

	// wait for cmd1
	cmd1Done := make(chan error, 1)

	go func() {
		cmd1Done <- cmd1.Wait()
		writer.Close()
	}()

	// wait for last command
	err := <-cmd2Done

	// stop infinite producers (tail -f)
	if cmd1.Process != nil {
		_ = cmd1.Process.Kill()
	}

	<-cmd1Done

	reader.Close()
	writer.Close()

	if err != nil {
		return "", ""
	}

	return "", ""
}

type LineWriter struct {
	w io.Writer
}

func (lw LineWriter) Write(p []byte) (int, error) {
	p = bytes.ReplaceAll(p, []byte("\n"), []byte("\r\n"))
	return lw.w.Write(p)
}

func HandleExternal(command string, args []string, background bool) (string, string) {
	found, _ := utils.ScanPath(os.Getenv("PATH"), command)

	if found == nil {
		return "", command + ": command not found\n"
	}

	cmd := exec.Command(*found, args...)

	return RunCommand(
		cmd,
		command+" "+strings.Join(args, " "),
		background,
	)
}
