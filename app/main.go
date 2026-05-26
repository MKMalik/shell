package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtinNames = map[Builtin]struct{}{
	// Echo: {},
	// Type: {},
}

var builtins = map[Builtin]func(string){
	Echo: handleEcho,
	Type: handleType,
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
		fmt.Println(arg + ": not found")
		return
	}

	fmt.Println(arg + " is a shell builtin")
}
func handleEcho(cmd string) {
	args := strings.Fields(cmd)[1:]
	fmt.Println(strings.Join(args, " "))
}
