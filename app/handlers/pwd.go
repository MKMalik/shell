package handlers

import (
	"fmt"
	"os"
)

func HandlePwd(cmd string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dir)
}
