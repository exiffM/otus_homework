package main

import (
	"fmt"
	"os"
)

var args []string

func init() {
	args = os.Args
}

func main() {
	if len(args) == 1 {
		fmt.Println("You have insert not enough arguments for utility usage.")
		fmt.Println("First insert path to utility environment variables.")
		fmt.Println("Second insert command for utility to execute and arguments for this command.")
		return
	}
	if env, err := ReadDir(args[1]); err != nil {
		fmt.Printf("Error occued: %v\n", err.Error())
	} else {
		RunCmd(args[2:], env)
	}
}
