package main

import (
	"fmt"
	"log"
	"quick-cmd/command"
	"quick-cmd/utils"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var supportCmd = []string{"bashHistory", "jumpDir", "git"}

func main() {
	cmd := utils.GetCmd()

	if len(cmd) == 0 {
		fmt.Println(`只支持命令:`, strings.Join(supportCmd, ", "))
		return
	}

	var err error
	switch cmd[0] {
	case "bashHistory":
		err = command.BashHistory()
	case "jumpDir":
		err = command.JumpDir()
	case "git":
		err = command.Git(cmd[1:])
	default:
		fmt.Println(`只支持命令:`, strings.Join(supportCmd, ", "))
		return
	}

	if err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}

}
