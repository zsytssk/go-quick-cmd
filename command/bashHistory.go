package command

import (
	"fmt"
	"quick-cmd/dbt"
	"quick-cmd/utils"
	"strings"
)

// BashHistoryCommand 实现了bash历史命令
type BashHistoryCommand struct {
	*BaseCommand
}

// NewBashHistoryCommand 创建新的bash历史命令
func NewBashHistoryCommand() (*BashHistoryCommand, error) {
	base, err := NewBaseCommand("bashHistory")
	if err != nil {
		return nil, err
	}
	return &BashHistoryCommand{BaseCommand: base}, nil
}

// Execute 执行bash历史命令
func (b *BashHistoryCommand) Execute() error {
	items, err := dbt.GetHistory(b.DB)
	if err != nil {
		return fmt.Errorf("failed to get history items: %w", err)
	}

	var fzfInput strings.Builder
	for _, item := range items {
		fzfInput.WriteString(fmt.Sprintf("%s [%d:%d]\n", item.Name, item.ID, item.Priority))
	}

	selected, err := utils.RunFZF(fzfInput.String())
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
	if err := dbt.UpdateHistoryPriority(b.DB, item); err != nil {
		return fmt.Errorf("failed to update priority: %w", err)
	}

	fmt.Print(item.Name)
	return nil
}
