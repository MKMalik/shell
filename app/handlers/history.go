package handlers

import (
	"fmt"
	"strings"
)

var HistoryList []string = make([]string, 0)

func HandleHistory(cmd string) string {
	result := make([]string, 0)
	for index, val := range HistoryList {
		result = append(result, fmt.Sprintf("    %v  %v", index+1, val))
	}
	return strings.Join(result, "\n")
}

func AppendHistory(cmd string) {
	HistoryList = append(HistoryList, cmd)
}

func reverse(s []string) []string {
	out := make([]string, len(s))

	for i := range s {
		out[len(s)-1-i] = s[i]
	}

	return out
}
