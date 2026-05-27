package handlers

import (
	"os"
)

func HandlePwd(cmd string) string {
	dir, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	return dir
}
