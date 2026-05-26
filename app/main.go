package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var builtinNames = map[Builtin]struct{}{}

var builtins = map[Builtin]func(string){
	Echo: handleEcho,
	Type: handleType,
	Exit: func(s string) {},
}

func isBuiltin(cmd string) bool {
	_, ok := builtinNames[Builtin(cmd)]
	return ok
}

type Builtin string

const (
	Exit Builtin = "exit"
	Echo Builtin = "echo"
	Type Builtin = "type"
)

func main() {
	for i := range builtins {
		builtinNames[i] = struct{}{}
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
	case string(Echo):
		builtins[Echo](command)
		return
	case string(Type):
		builtins[Type](command)
		return
	default:
		// check if exists in PATH as executable
		found, _ := scanPath(os.Getenv("PATH"), cmd[0])
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

func handleType(cmd string) {
	fields := strings.Fields(cmd)
	if len(fields) < 2 {
		return
	}
	arg := fields[1]
	if !isBuiltin(arg) {
		paths := os.Getenv("PATH")
		_, full := scanPath(paths, arg)
		if full != nil {
			fmt.Println(arg + " is " + *full)
			return
		}
		fmt.Println(arg + ": not found")
		return
	}

	fmt.Println(arg + " is a shell builtin")
}

func scanPath(paths, arg string) (*string, *string) {
	entries := strings.Split(paths, ":")
	for i := range entries {
		dirEntries, err := os.ReadDir(entries[i])
		if err != nil {
			// panic(err)
			continue
		}
		for _, dirEnt := range dirEntries {
			if dirEnt.IsDir() {
				scanPath(entries[i]+"/"+dirEnt.Name(), arg)
			} else {
				info, err := dirEnt.Info()
				if err != nil {
					log.Printf("Error getting info for %s: %v", dirEnt.Name(), err)
					continue
				}
				var isExec bool = info.Mode().Perm()&0100 != 0
				if isExec && arg == dirEnt.Name() {
					found := dirEnt.Name()
					full := entries[i] + "/" + found
					return &found, &full
				}
			}
		}
	}
	return nil, nil
}

func handleEcho(cmd string) {
	args := strings.Fields(cmd)[1:]
	fmt.Println(strings.Join(args, " "))
}
