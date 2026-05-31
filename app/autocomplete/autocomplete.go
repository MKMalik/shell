package autocomplete

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/handlers"
)

var first bool = true

func HandleCommandAutocomplete(input []byte) []byte {
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

func HandleFileAutocomplete(input []byte) []byte {
	line := string(input)

	fields := strings.Split(line, " ")
	if len(fields) < 2 {
		return input
	}

	prefix := fields[len(fields)-1]

	matches := FindFilesAndDirs(prefix)

	switch len(matches) {
	case 0:
		os.Stdout.WriteString("\a")
		return input

	case 1:
		completed := matches[0]

		info, err := os.Stat(completed)
		if err == nil && info.IsDir() {
			completed += "/"
		} else {
			completed += " "
		}

		fields[len(fields)-1] = completed

		newInput := strings.Join(fields, " ")

		os.Stdout.WriteString("\r\033[2K$ ")
		os.Stdout.WriteString(newInput)

		return []byte(newInput)

	default:
		lcp := LongestCommonPrefix(matches)

		if len(lcp) > len(prefix) {
			fields[len(fields)-1] = lcp

			newInput := strings.Join(fields, " ")

			os.Stdout.WriteString("\r\033[2K$ ")
			os.Stdout.WriteString(newInput)

			return []byte(newInput)
		}

		if first {
			first = false
			os.Stdout.WriteString("\a")
			return input
		}
		first = true

		os.Stdout.WriteString("\r\n")
		display := make([]string, 0, len(matches))
		for _, m := range matches {
			display = append(display, DisplayName(m))
		}
		os.Stdout.WriteString(strings.Join(display, "  "))
		os.Stdout.WriteString("\r\n")
		os.Stdout.WriteString("$ ")
		os.Stdout.WriteString(string(input))

		return input
	}
}

func FindFilesAndDirs(prefix string) []string {
	dir := "."
	base := prefix

	if before, ok := strings.CutSuffix(prefix, "/"); ok {
		dir = before
		base = ""
	} else if strings.Contains(prefix, "/") {
		dir = filepath.Dir(prefix)
		base = filepath.Base(prefix)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var matches []string

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), base) {
			continue
		}

		full := filepath.Join(dir, entry.Name())

		matches = append(matches, full)
	}

	sort.Strings(matches)
	return matches
}

func DisplayName(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return filepath.Base(path)
	}

	if info.IsDir() {
		return filepath.Base(path) + "/"
	}

	return filepath.Base(path)
}
