package autocomplete

import (
	"os"
	"sort"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
)

var first bool = true

func HandleAutocomplete(input []byte) []byte {
	seen := make(map[string]struct{})
	var matches []string

	for val := range handlers.Builtins {
		name := string(val)

		if !strings.HasPrefix(name, string(input)) {
			continue
		}

		seen[name] = struct{}{}
		matches = append(matches, name)
	}

	for _, name := range FindExecutables(string(input)) {
		if _, ok := seen[name]; ok {
			continue
		}

		seen[name] = struct{}{}
		matches = append(matches, name)
	}

	sort.Strings(matches)

	switch len(matches) {
	case 0:
		os.Stdout.WriteString("\x07")
		return input
	case 1:
		os.Stdout.WriteString("\r\033[2K$ ")
		os.Stdout.WriteString(matches[0])
		os.Stdout.WriteString(" ")
		return []byte(matches[0] + " ")
	default:
		lcp := LongestCommonPrefix(matches)

		if len(lcp) > len(input) {
			suffix := lcp[len(input):]

			os.Stdout.WriteString(suffix)

			return []byte(lcp)
		}

		if first {
			first = false
			os.Stdout.WriteString("\a")
			return input
		}
		first = true

		os.Stdout.WriteString("\r\n")
		os.Stdout.WriteString(strings.Join(matches, "  "))
		os.Stdout.WriteString("\r\n")
		os.Stdout.WriteString("$ ")
		os.Stdout.WriteString(string(input))

		return input
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

func LongestCommonPrefix(matches []string) string {
	if len(matches) == 0 {
		return ""
	}

	prefix := matches[0]

	for _, s := range matches[1:] {
		for !strings.HasPrefix(s, prefix) {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}

	return prefix
}
