package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, err := scanner.ReadString('\n')
		if err != nil {
			panic(err)
		}
		processCmd(cmd)
	}
}

func processCmd(cmd string) {
	fmt.Println(cmd[:len(cmd)-1] + ": command not found")
}
