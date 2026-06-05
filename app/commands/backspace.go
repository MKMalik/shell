package commands

import "os"

func HandleBackspace(input []byte) []byte {
	if len(input) == 0 {
		return input
	}

	input = input[:len(input)-1]
	os.Stdout.WriteString("\b \b")

	return input
}
