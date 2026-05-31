package main

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/commands"
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
				s := string(input)

				if strings.ContainsRune(s, ' ') {
					input = autocomplete.HandleFileAutocomplete(input)
				} else {
					input = autocomplete.HandleCommandAutocomplete(input)
				}

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

		redirect := utils.ParseRedirect(cmd)

		if redirect.Valid {
			stdout, stderr := processCmd(redirect.Cmd)

			if redirect.FD == 1 {
				_ = utils.RedirectToFile(stdout, redirect.File, redirect.Append)

				if stderr != "" {
					utils.WriteOutput(stderr)
				}
			} else {
				// FD == 2
				_ = utils.RedirectToFile(stderr, redirect.File, redirect.Append)

				if stdout != "" {
					utils.WriteOutput(stdout)
				}
			}
			continue
		}

		stdout, stderr := processCmd(cmd)
		if stdout != "" {
			utils.WriteOutput(stdout)
		} else {
			utils.WriteOutput(stderr)
		}
	}
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
		return HandleExternal(fields[0], args)
	}
}

func HandleExternal(command string, args []string) (string, string) {
	found, _ := utils.ScanPath(os.Getenv("PATH"), command)

	if found != nil {
		stdout, stderr, _ := commands.RunCommand(*found, args)
		return stdout, stderr
	}

	return "", command + ": command not found\n"
}
