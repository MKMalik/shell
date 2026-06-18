package handlers

import (
	"fmt"
	"strings"
)

var declareValues map[string]string = map[string]string{}

func HandleDeclare(cmd string) string {
	fields := strings.Fields(cmd)
	declared := fields[1]
	if declared == "-p" {

		declareVar := fields[2]
		if _, ok := declareValues[declareVar]; ok {
			return fmt.Sprintf("declare -- %v=\"%v\"", declareVar, declareValues[declareVar])
		}
		return fmt.Sprintf("declare: %v: not found", declareVar)

	} else if strings.Contains(declared, "=") {

		splits := strings.SplitN(declared, "=", 2)
		key, val := splits[0], splits[1]

		if !isValidIdentifier(key) {
			return fmt.Sprintf("declare: `%v': not a valid identifier", declared)
		}

		declareValues[key] = val

	}
	return ""
}

func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}

	if !((s[0] >= 'a' && s[0] <= 'z') ||
		(s[0] >= 'A' && s[0] <= 'Z') ||
		s[0] == '_') {
		return false
	}

	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_') {
			return false
		}
	}

	return true
}
