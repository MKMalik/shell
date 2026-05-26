package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	handlers "github.com/codecrafters-io/shell-starter-go/app/handlers"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func main() {
	for i := range handlers.Builtins {
		handlers.BuiltinNames[i] = struct{}{}
	}
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, err := scanner.ReadString('\n')
		if err != nil {
			panic(err)
		}
		cmd = strings.TrimSpace(cmd)
		if cmd == "exit" {
			break
		}
		processCmd(cmd)
	}
}

func processCmd(command string) {
	cmd := strings.Split(command, " ")
	args := cmd[1:]
	switch cmd[0] {
	case string(handlers.Echo):
		handlers.Builtins[handlers.Echo](command)
		return
	case string(handlers.Type):
		handlers.Builtins[handlers.Type](command)
		return
	case string(handlers.Pwd):
		handlers.Builtins[handlers.Pwd](command)
	case string(handlers.Cd):
		handlers.Builtins[handlers.Cd](args[0])
	default:
		// check if exists in PATH as executable
		found, _ := utils.ScanPath(os.Getenv("PATH"), cmd[0])
		// if exists and executable then execute passing args if any
		if found != nil {
			run := exec.Command(*found, args...)

			stdout, err := run.StdoutPipe()
			if err != nil {
				return
			}

			stderr, err := run.StderrPipe()
			if err != nil {
				return
			}

			if err := run.Start(); err != nil {
				return
			}

			stdoutBytes, _ := io.ReadAll(stdout)
			stderrBytes, _ := io.ReadAll(stderr)

			if err := run.Wait(); err != nil {
				fmt.Println(err)
			}

			fmt.Print(string(stdoutBytes))
			fmt.Print(string(stderrBytes))
			return
		}
		// if not: print command not found
		fmt.Println(command + ": command not found")
		return
	}
}
