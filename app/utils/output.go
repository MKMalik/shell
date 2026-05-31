package utils

import (
	"os"
	"strings"
)

func WriteOutput(output string) {
	output = strings.TrimRight(output, "\n")

	if output == "" {
		return
	}

	for line := range strings.SplitSeq(output, "\n") {
		os.Stdout.WriteString("\r")
		os.Stdout.WriteString(line)
		os.Stdout.WriteString("\r\n")
	}
}
