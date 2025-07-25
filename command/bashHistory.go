package command

import (
	"fmt"
	"quick-cmd/dbt"
	"quick-cmd/utils"
	"sort"
	"strings"
)

type HistoryItem struct {
	Item
	Content string
}

func (HistoryItem) TableName() string {
	return "history"
}

func BashHistory() (err error) {
	db, err := getDb()
	if err != nil {
		return
	}
	dm, err := dbt.NewModel(db, &HistoryItem{})
	if err != nil {
		return fmt.Errorf("failed to get history items: %w", err)
	}
	items, err := GetHistory(dm)
	if err != nil {
		return fmt.Errorf("failed to get history items: %w", err)
	}

	var fzfInput strings.Builder
	for _, item := range items {
		if item.Hide {
			continue
		}
		fzfInput.WriteString(fmt.Sprintf("%s [%d:%d]\n", item.Name, item.ID, item.Priority))
	}

	action, selected, err := utils.RunFZF(fzfInput.String())
	if err != nil {
		if utils.IsCanceled(err) {
			return nil
		}
		return fmt.Errorf("failed to run fzf: %w", err)
	}

	if selected == "" {
		return nil
	}

	item, found := utils.ArrFind(items, func(item HistoryItem, _ int) bool {
		return selected == fmt.Sprintf("%s [%d:%d]", item.Name, item.ID, item.Priority)
	})

	if !found {
		return fmt.Errorf("item not found: %s", selected)
	}

	if action == utils.ActionDelete {
		item.Hide = true
		if err := dm.Save(item).Error; err != nil {
			return fmt.Errorf("failed to save item: %w", err)
		}
		return BashHistory()
	}

	item.Priority = item.Priority + 1
	if err := dm.Save(item).Error; err != nil {
		return fmt.Errorf("failed to save item: %w", err)
	}
	if len(item.Content) > 0 {
		fmt.Print(item.Content)
		return
	}

	fmt.Print(item.Name)
	return
}

func GetHistory(dm *dbt.Model) (items []HistoryItem, err error) {
	lineMap, err := utils.ReadFileLines("~/.bash_history")

	if err != nil {
		return
	}

	for key := range lineMap {
		if strings.HasPrefix(key, "cd") && !strings.Contains(key, "&&") {
			delete(lineMap, key)
			continue
		}
	}
	var count int64
	err = dm.Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		for key := range lineMap {
			if strings.HasPrefix(key, "cd") && !strings.Contains(key, "&&") {
				dm.Save(DirItem{Item{-1, key, lineMap[key], false}})
				continue
			}
		}
	}
	err = dm.Order("ORDER BY priority DESC").Find(&items).Error
	if err != nil {
		return
	}
	for key, count := range lineMap {
		index := utils.ArrFindIndex(items, func(item HistoryItem, _ int) bool {
			return item.Name == key
		})
		if index != -1 {
			continue
		}
		// fmt.Println("test:>", key, count)
		item := HistoryItem{Item{-1, key, count, false}, ""}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	// index := utils.FindItemIndex(items, func(item Item, _ int) bool {
	// 	return item.Name == "rm ./go-test"
	// })
	// fmt.Println("test:>", items)

	return
}
