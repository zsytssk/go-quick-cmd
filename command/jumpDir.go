package command

import (
	"fmt"
	"io"
	"quick-cmd/dbt"
	"quick-cmd/utils"
	"strings"
)

// JumpDirCommand 实现了目录跳转命令
type JumpDirCommand struct {
	*BaseCommand
}

// NewJumpDirCommand 创建新的目录跳转命令
func NewJumpDirCommand() (*JumpDirCommand, error) {
	base, err := NewBaseCommand("jumpDir")
	if err != nil {
		return nil, err
	}
	return &JumpDirCommand{BaseCommand: base}, nil
}

// Execute 执行目录跳转命令
func (j *JumpDirCommand) Execute() error {
	config, err := utils.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	items, err := dbt.GetDir(j.DB)
	if err != nil {
		return fmt.Errorf("failed to get dir items: %w", err)
	}

	cmdStr := buildFindStr(config)
	reader, writer := io.Pipe()
	defer reader.Close()

	go func() {
		defer writer.Close()
		for _, item := range items {
			fmt.Fprintf(writer, "%s [%d:%d]\n", item.Name, item.ID, item.Priority)
		}
		utils.RunCMDInSteam(cmdStr, func(line string) {
			index := utils.ArrFindIndex(items, func(item dbt.Item, _ int) bool {
				return item.Name == line
			})
			if index != -1 {
				return
			}
			item := dbt.Item{ID: -1, Name: line, Priority: 0}
			items = append(items, item)
			fmt.Fprintf(writer, "%s [%d:%d]\n", item.Name, item.ID, item.Priority)
		})
	}()

	selected, err := utils.RunFZFStream(reader)
	if err != nil {
		if utils.IsCanceled(err) {
			return nil
		}
		return fmt.Errorf("failed to run fzf: %w", err)
	}

	if selected == "" {
		return nil
	}

	index := utils.ArrFindIndex(items, func(item dbt.Item, _ int) bool {
		return selected == fmt.Sprintf("%s [%d:%d]", item.Name, item.ID, item.Priority)
	})

	if index == -1 {
		return fmt.Errorf("item not found: %s", selected)
	}

	item := items[index]
	if err := dbt.UpdateDirPriority(j.DB, item); err != nil {
		return fmt.Errorf("failed to update priority: %w", err)
	}

	fmt.Print(`cd `, item.Name)
	return nil
}

// buildFindStr 构建find命令字符串
func buildFindStr(config utils.Config) string {
	var cmdInput strings.Builder
	for _, item := range config.Folders {
		ignores := append(config.Ignores, item.Ignores...)
		ignoreStr := utils.ArrJoin(ignores, func(item string, index int) string {
			if index == 0 {
				return fmt.Sprintf(` -path %s`, item)
			}
			return fmt.Sprintf(` -o -path %s`, item)
		})
		cmdInput.WriteString(fmt.Sprintf("find %s -maxdepth %d -type d \\( %s \\)  -prune -o -print\n", item.Folder, item.Depth, ignoreStr))
	}
	return cmdInput.String()
}
