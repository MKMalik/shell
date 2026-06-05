package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"golang.org/x/term"
)

func main() {
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 1)


	for {
		cmd, ok := commands.ReadCommand(buf)
		if !ok {
			return
		}

		commands.ExecuteCommand(cmd)
	}
}
