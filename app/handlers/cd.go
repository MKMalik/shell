package handlers

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleCd(cmd string) string {
	dir := strings.Split(cmd, " ")[1]
	if dir == "~" {
		return cdHomeDir()
	} else if string(dir[0]) == "/" {
		return cdAbsoluteDir(dir)
	} else {
		return cdRelativeDir(dir)
	}
}

func cdRelativeDir(dir string) string {
	pwd, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	if exists := utils.IsDirExists(pwd + "/" + dir); !exists {
		return "cd: /" + dir + ": No such file or directory" + "\n"
	}
	os.Chdir(dir)
	return ""
}

func cdAbsoluteDir(dir string) string {
	if exists := utils.IsDirExists(dir); !exists {
		return "cd: " + dir + ": No such file or directory" + "\n"
	}
	os.Chdir(dir)
	return ""
}

func cdHomeDir() string {
	home := os.Getenv("HOME")
	err := os.Chdir(home)
	if err != nil {
		return err.Error()
	}
	os.Chdir(home)
	return ""
}
