package handlers

import (
	"fmt"
	"strings"

	"github.com/google/shlex"
)

func HandleEcho(cmd string) {
	args, err := shlex.Split(cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.Join(args[1:], " "))
}
