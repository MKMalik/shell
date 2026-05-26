package handlers

import (
	"fmt"
	"os"
	"os/user"
)

func HandleCd(dir string) {
	if dir == "~" {
		usr, err := user.Current()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = os.Chdir(usr.HomeDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	os.Chdir(dir)
}
