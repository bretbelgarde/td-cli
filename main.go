package main

import (
	"fmt"
	"os"

	"bretbelgarde.com/td-cli/cmd"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Expected one of the following: 'add', 'list', 'delete', 'update', or 'complete'")
		os.Exit(1)
	}

	cmd.Execute(os.Args)
	os.Exit(0)
}
