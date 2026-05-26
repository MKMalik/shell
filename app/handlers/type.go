package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleType(cmd string) {
	fields := strings.Fields(cmd)
	if len(fields) < 2 {
		return
	}
	arg := fields[1]
	if !IsBuiltin(arg) {
		paths := os.Getenv("PATH")
		_, full := utils.ScanPath(paths, arg)
		if full != nil {
			fmt.Println(arg + " is " + *full)
			return
		}
		fmt.Println(arg + ": not found")
		return
	}

	fmt.Println(arg + " is a shell builtin")
}
