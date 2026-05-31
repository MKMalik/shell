package handlers

import (
	"fmt"
	"strings"
)

func HandleComplete(cmd string) string {
	split := strings.Split(cmd, "-p")
	if len(split) < 2 {
		return ""
	}
	arg := strings.TrimSpace(split[1])
	return fmt.Sprintf("complete: %v: no completion specification", arg)
}
