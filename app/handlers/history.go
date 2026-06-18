package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var HistoryList []string = make([]string, 0)
var currentHistoryIndex = len(HistoryList)
var historyWriteIndex = 0

func HandleHistory(cmd string) string {
	args := strings.Fields(cmd)

	limit := len(HistoryList)

	if len(args) > 1 && args[1] == "-r" {
		if len(args) < 3 {
			return ""
		}

		file, err := os.ReadFile(args[2])
		if err != nil {
			return ""
		}

		for val := range strings.SplitSeq(string(file), "\n") {
			if val == "" {
				continue
			}
			HistoryList = append(HistoryList, val)
		}

		currentHistoryIndex = len(HistoryList)

		return ""
	}

	if len(args) > 1 && args[1] == "-w" {
		os.WriteFile(
			args[2],
			[]byte(strings.Join(HistoryList, "\n")+"\n"),
			0644,
		)
		return ""
	}

	if len(args) > 1 && args[1] == "-a" {
		if len(args) < 3 {
			return ""
		}

		f, err := os.OpenFile(
			args[2],
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0644,
		)
		if err != nil {
			return ""
		}
		defer f.Close()

		if historyWriteIndex < len(HistoryList) {
			_, _ = f.WriteString(
				strings.Join(HistoryList[historyWriteIndex:], "\n") + "\n",
			)

			historyWriteIndex = len(HistoryList)
		}

		return ""
	}

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
