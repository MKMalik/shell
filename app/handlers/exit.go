package handlers

import "os"

func HandleExit() {
	WriteHistoryToFile(os.Getenv("HISTFILE"))
	os.Exit(0)
}
