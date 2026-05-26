package handlers

import (
	"fmt"
	"os"
	"os/user"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
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

	// if absolute path
	if string(dir[0]) == "/" {
		if exists := utils.IsDirExists(dir); !exists {
			fmt.Println("cd: " + dir + ": No such file or directory")
			return
		}
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
		if exists := utils.IsDirExists(pwd + "/" + dir); !exists {
			fmt.Println("cd: " + dir + ": No such file or directory")
			return
		}
	}
	os.Chdir(dir)
}
