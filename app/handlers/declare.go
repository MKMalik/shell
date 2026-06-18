package handlers

import (
	"fmt"
	"strings"
)

func HandleDeclare(cmd string) string {
	fields := strings.Fields(cmd)
	// declareType := fields[1]
	declareVar := fields[2]
	return fmt.Sprintf("declare: %v: not found", declareVar)
}
