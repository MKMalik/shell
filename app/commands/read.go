package commands

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
)

func ReadCommand(buf []byte) (string, bool) {
	os.Stdout.WriteString("\r\033[2K$ ")

	var input []byte

	for {
		_, _ = os.Stdin.Read(buf)

		switch buf[0] {
		case 13, 10:
			os.Stdout.WriteString("\r\n")
			handlers.AppendHistory(string(input))
			return string(input), true

		case 9:
			input = HandleTab(input)

		case 127, 8:
			input = HandleBackspace(input)

		case 3:
			os.Exit(1)

		default:
			input = append(input, buf[0])
			os.Stdout.Write(buf)
		}
	}
}
