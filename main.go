package main

import (
	"fmt"
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

	switch *cmd {
	case "bashHistory":
		command.HashHistory()
	case "jumpDir":
		command.JumpDir()
	default:
		command.JumpDir()
	}
}
