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
		splits := strings.Split(declared, "=")
		if len(splits) < 2 {
			return ""
		}
		key, val := splits[0], splits[1]
		declareValues[key] = val
	}
	return ""
}
