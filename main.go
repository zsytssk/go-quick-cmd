package main

import (
	"fmt"
	"go-sqlite-test/command"
	"go-sqlite-test/utils"

	_ "github.com/mattn/go-sqlite3"
)

var supportCmd = []string{"bashHistory", "JumpDir"}

func main() {

	cmd := utils.GetCmd()
	if cmd != nil && !utils.ArrContains(supportCmd, *cmd) {
		fmt.Printf("不支持命令 %v \n", *cmd)
		return
	}

	switch *cmd {
	case "bashHistory":
		command.HashHistory()
	case "JumpDir":
		command.JumpDir()
	default:
		command.JumpDir()
	}
}
