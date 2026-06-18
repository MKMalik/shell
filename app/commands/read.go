package commands

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
)

func ReadCommand(buf []byte) (string, bool) {
	// read history from HISTFILE env variable and append into HistoryList memory variable
	LoadHistory()
	os.Stdout.WriteString("\r\033[2K$ ")
	var input []byte
	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}
		switch {
		case buf[0] == 13 || buf[0] == 10:
			os.Stdout.WriteString("\r\n")
			handlers.AppendHistory(string(input))
			return string(input), true
		case buf[0] == 9:
			input = HandleTab(input)
		case buf[0] == 127 || buf[0] == 8:
			input = HandleBackspace(input)
		case buf[0] == 3:
			handlers.WriteHistoryToFile(os.Getenv("HISTFILE"))
			os.Exit(1)
		case buf[0] == 27:
			// read 2 more bytes to complete the escape sequence
			seq := make([]byte, 2)
			os.Stdin.Read(seq)
			if seq[0] == 91 {
				switch seq[1] {
				case 65:
					input = handlers.HistoryUp()
					os.Stdout.WriteString("\r\033[2K$ ")
					os.Stdout.Write(input)
				case 66:
					input = handlers.HistoryDown()
					os.Stdout.WriteString("\r\033[2K$ ")
					os.Stdout.Write(input)
				}
			}
		default:
			input = append(input, buf[0])
			os.Stdout.Write(buf[:1])
		}
	}
}

func LoadHistory() {
	file, err := os.ReadFile(os.Getenv("HISTFILE"))
	if err != nil {
		return
	}

	for val := range strings.SplitSeq(string(file), "\n") {
		if val == "" {
			continue
		}
		handlers.HistoryList = append(handlers.HistoryList, val)
	}

}
