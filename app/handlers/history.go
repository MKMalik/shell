package handlers

import (
	"fmt"
	"strconv"
	"strings"
)

var HistoryList []string = make([]string, 0)
var currentHistoryIndex = len(HistoryList)

func HandleHistory(cmd string) string {
	args := strings.Fields(cmd)

	limit := len(HistoryList)

	if len(args) == 2 {
		if n, err := strconv.Atoi(args[1]); err == nil {
			limit = n
		}
	}

	start := max(len(HistoryList)-limit, 0)

	var history []string

	for i := start; i < len(HistoryList); i++ {
		history = append(history, fmt.Sprintf(
			"    %d  %s",
			i+1,
			HistoryList[i],
		))
	}

	return strings.Join(history, "\n")
}

func HistoryUp() []byte {
	if len(HistoryList) == 0 {
		return []byte{}
	}

	currentHistoryIndex--

	if currentHistoryIndex < 0 {
		currentHistoryIndex = 0
	}

	return []byte(HistoryList[currentHistoryIndex])
}

func HistoryDown() []byte {
	if len(HistoryList) == 0 {
		return []byte{}
	}

	currentHistoryIndex++

	if currentHistoryIndex >= len(HistoryList) {
		currentHistoryIndex = len(HistoryList)
		return []byte("")
	}

	return []byte(HistoryList[currentHistoryIndex])
}

func LastN(list []string, n int) []string {
	if n <= 0 {
		return []string{}
	}

	if n > len(list) {
		n = len(list)
	}

	return list[len(list)-n:]
}

func AppendHistory(cmd string) {
	HistoryList = append(HistoryList, cmd)
	currentHistoryIndex = len(HistoryList)
}
