package main

import (
	"fmt"
	"log"
	"quick-cmd/command"
	"quick-cmd/utils"
	"slices"

	_ "github.com/mattn/go-sqlite3"
)

var supportCmd = []string{"bashHistory", "jumpDir"}

func main() {
	cmd := utils.GetCmd()
	if cmd == nil {
		fmt.Println(`请输入执行命令 "bashHistory" | "jumpDir"`)
		return
	}
	if !slices.Contains(supportCmd, *cmd) {
		fmt.Println(`只支持命令："bashHistory" | "jumpDir"`, *cmd)
		return
	}

	var err error

	switch *cmd {
	case "bashHistory":
		err = command.BashHistory()
	default:
		err = command.JumpDir()
	}

	if err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}

}
