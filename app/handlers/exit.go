package handlers

import "os"

func HandleExit() {
	AppendHistoryToFile(os.Getenv("HISTFILE"))
	os.Exit(0)
}
