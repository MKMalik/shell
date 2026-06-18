package handlers

import "strings"

var HistoryList []string = make([]string, 0)

func HandleHistory(cmd string) string {
	return strings.Join(reverse(HistoryList), "\n")
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
