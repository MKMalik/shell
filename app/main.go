package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	fmt.Print("$ ")
	scanner := bufio.NewReader(os.Stdin)
	cmd, err := scanner.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd[:len(cmd) -1 ] + ": command not found")
}
