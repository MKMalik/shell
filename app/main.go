package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func processCmd(cmd string) {
	switch strings.Split(cmd, " ")[0] {
	case string(Echo):
		builtins[Echo](cmd)
	case string(Type):
		builtins[Type](cmd)
	default:
		fmt.Println(cmd + ": command not found")
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
		found := scanPath(paths, arg)
		if found != nil {
			fmt.Println(arg + " is " + *found)
			return
		}
		fmt.Println(arg + ": not found")
		return
	}

	fmt.Println(arg + " is a shell builtin")
}

func scanPath(paths, arg string) *string {
	entries := strings.Split(paths, ":")
	for i := range entries {
		dirEntries, err := os.ReadDir(entries[i])
		if err != nil {
			panic(err)
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
					found := entries[i] + "/" + dirEnt.Name()
					return &found
				}
			}
		}
	}
	return nil
}
func handleEcho(cmd string) {
	args := strings.Fields(cmd)[1:]
	fmt.Println(strings.Join(args, " "))
}
