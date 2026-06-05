package commands

import (
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/app/handlers"
)

func HandleTab(input []byte) []byte {
	s := string(input)

	fields := strings.Fields(s)

	if len(fields) > 0 && handlers.GetComplete(fields[0]) != nil {
		out, result := handlers.RunCompleter(s)

		switch result {
		case handlers.Completed:
			return []byte(out)

		case handlers.Handled:
			return input

		case handlers.NoCompletion:
		}
	}

	if strings.ContainsRune(s, ' ') {
		return autocomplete.HandleFileAutocomplete(input)
	}

	return autocomplete.HandleCommandAutocomplete(input)
}
