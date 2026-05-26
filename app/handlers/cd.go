package handlers

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleCd(dir string) {
	if dir == "~" {
		cdHomeDir()
	} else if string(dir[0]) == "/" {
		cdAbsoluteDir(dir)
	} else {
		cdRelativeDir(dir)
	}
}

func cdRelativeDir(dir string) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	if exists := utils.IsDirExists(pwd + "/" + dir); !exists {
		fmt.Println("cd: /" + dir + ": No such file or directory")
		return
	}
	os.Chdir(dir)
}

func cdAbsoluteDir(dir string) {
	if exists := utils.IsDirExists(dir); !exists {
		fmt.Println("cd: " + dir + ": No such file or directory")
	}
	os.Chdir(dir)
}

func cdHomeDir() {
	home := os.Getenv("HOME")
	err := os.Chdir(home)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Chdir(home)
}
