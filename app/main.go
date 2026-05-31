package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"github.com/google/shlex"
	"golang.org/x/term"
)

func main() {
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 1)

	for {
		os.Stdout.Write([]byte("\r\033[2K$ "))

		var input []byte

		for {
			os.Stdin.Read(buf)

			if buf[0] == 13 || buf[0] == 10 {
				os.Stdout.Write([]byte("\r\n"))
				break
			}

			if buf[0] == 9 {
				input = handleTab(input)
				continue
			}

			if buf[0] == 127 || buf[0] == 8 {
				if len(input) > 0 {
					input = input[:len(input)-1]
					os.Stdout.Write([]byte("\b \b"))
				}
				continue
			}

			input = append(input, buf[0])
			os.Stdout.Write(buf)
		}

		cmd := string(input)

		if strings.TrimSpace(cmd) == "exit" {
			return
		}
		redirect := parseRedirect(cmd)

		if redirect.Valid {
			// fmt.Printf("CMD=[%s]\n", redirect.Cmd)
			stdout, stderr := processCmd(redirect.Cmd)

			if redirect.FD == 1 {
				_ = redirectToFile(stdout, redirect.File, redirect.Append)

				if stderr != "" {
					writeOutput(stderr)
				}
			} else {
				_ = redirectToFile(stderr, redirect.File, redirect.Append)

				if stdout != "" {
					writeOutput(stdout)
				}
			}
			continue
		}

		stdout, stderr := processCmd(cmd)
		if stdout != "" {
			writeOutput(stdout)
		} else {
			writeOutput(stderr)
		}
	}
}

func handleTab(input []byte) []byte {
	for val := range handlers.Builtins {
		if strings.HasPrefix(string(val), string(input)) {
			os.Stdout.WriteString("\r\033[2K$ ")
			os.Stdout.WriteString(string(val + " "))
			return []byte(val + " ")
		}
	}
	os.Stdout.WriteString("\x07")
	return input
}

func writeOutput(output string) {
	output = strings.TrimRight(output, "\n")

	if output == "" {
		return
	}

	for _, line := range strings.Split(output, "\n") {
		os.Stdout.WriteString("\r")
		os.Stdout.WriteString(line)
		os.Stdout.WriteString("\r\n")
	}
}

func redirectToFile(content, file string, append bool) error {
	f, err := openFile(file, append)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func openFile(name string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return os.Create(name)
}

func processCmd(command string) (string, string) {

	fields, _ := shlex.Split(command)

	if len(fields) == 0 {
		return "", ""
	}

	args := fields[1:]

	switch fields[0] {
	case string(handlers.Echo):
		return handlers.Builtins[handlers.Echo](command), ""

	case string(handlers.Type):
		return handlers.Builtins[handlers.Type](command), ""

	case string(handlers.Pwd):
		return handlers.Builtins[handlers.Pwd](command), ""

	case string(handlers.Cd):
		return handlers.Builtins[handlers.Cd](args[0]), ""

	default:
		found, _ := utils.ScanPath(os.Getenv("PATH"), fields[0])

		if found != nil {
			stdout, stderr, _ := runCommand(*found, args)
			return stdout, stderr
		}

		return "", command + ": command not found\n"
	}
}

func runCommand(found string, args []string) (string, string, error) {
	run := exec.Command(found, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	run.Stdout = &stdout
	run.Stderr = &stderr

	err := run.Run()

	return stdout.String(), stderr.String(), err
}

type Redirect struct {
	Cmd    string
	File   string
	FD     int // 1=stdout, 2=stderr
	Append bool
	Valid  bool
}

func parseRedirect(cmd string) Redirect {
	ops := []struct {
		token  string
		fd     int
		append bool
	}{
		{"2>>", 2, true},
		{"1>>", 1, true},
		{">>", 1, true},
		{"2>", 2, false},
		{"1>", 1, false},
		{">", 1, false},
	}

	for _, op := range ops {
		parts := strings.SplitN(cmd, op.token, 2)

		if len(parts) != 2 {
			continue
		}

		return Redirect{
			Cmd:    strings.TrimSpace(parts[0]),
			File:   strings.TrimSpace(parts[1]),
			FD:     op.fd,
			Append: op.append,
			Valid:  true,
		}
	}

	return Redirect{
		Cmd:   cmd,
		Valid: false,
	}
}
