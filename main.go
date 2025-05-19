package main

import (
	"fmt"
	"log"
	"quick-cmd/command"
	"quick-cmd/utils"
	"slices"

	// 新增pty支持
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

	var cmdImpl command.Command
	var err error

	switch *cmd {
	case "bashHistory":
		cmdImpl, err = command.NewBashHistoryCommand()
	case "jumpDir":
		cmdImpl, err = command.NewJumpDirCommand()
	default:
		cmdImpl, err = command.NewJumpDirCommand()
	}

	if err != nil {
		log.Fatalf("Failed to create command: %v", err)
	}

	if err := cmdImpl.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}
