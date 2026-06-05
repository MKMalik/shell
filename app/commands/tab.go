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
		if out, ok := handlers.RunCompleter(s); ok {
			return []byte(out)
		}
	}

	if strings.ContainsRune(s, ' ') {
		return autocomplete.HandleFileAutocomplete(input)
	}

	return autocomplete.HandleCommandAutocomplete(input)
}
