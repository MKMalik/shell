package handlers

import "os"

func HandleCd(dir string) {
	os.Chdir(dir)
}
