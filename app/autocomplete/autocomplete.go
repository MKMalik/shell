package autocomplete

import (
	"os"
	"sort"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
)

func HandleAutocomplete(input []byte) []byte {
	var matches []string
	for val := range handlers.Builtins {
		if strings.HasPrefix(string(val), string(input)) {
			matches = append(matches, string(val))
		}
	}

	matches = append(matches, FindExecutables(string(input))...)

	switch len(matches) {
	case 0:
		os.Stdout.WriteString("\x07")
		return input
	case 1:
		os.Stdout.WriteString("\r\033[2K$ ")
		os.Stdout.WriteString(string(matches[0] + " "))
		return []byte(matches[0] + " ")
	default:
		os.Stdout.WriteString("\r\033[2K$ ")
		os.Stdout.WriteString(matches[0] + " ")
		return []byte(matches[0] + " ")
	}
}

func FindExecutables(prefix string) []string {
	seen := make(map[string]struct{})
	var matches []string

	for dir := range strings.SplitSeq(os.Getenv("PATH"), ":") {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()

			if !strings.HasPrefix(name, prefix) {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.Mode()&0111 == 0 {
				continue
			}

			if _, ok := seen[name]; ok {
				continue
			}

			seen[name] = struct{}{}
			matches = append(matches, name)
		}
	}

	sort.Strings(matches)
	return matches
}
