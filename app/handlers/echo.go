package handlers

import (
	"fmt"
	"strings"
)

func HandleEcho(cmd string) {
	args := strings.Fields(cmd)[1:]
	fmt.Println(strings.Join(args, " "))
}
