package command

import (
	"fmt"
	"io"
	"quick-cmd/dbt"
	"quick-cmd/utils"
	"sort"
	"strings"
)

type DirItem struct {
	Item
}

func (DirItem) TableName() string {
	return "dir"
}

func JumpDir() (err error) {
	db, err := getDb()
	if err != nil {
		return
	}
	dm := dbt.NewModel(db, &DirItem{})
	list, err := getDirHistory(dm)
	if err != nil {
		return
	}

	config, err := utils.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	cmdStr := buildFindStr(config)
	reader, writer := io.Pipe()
	defer reader.Close()

	go func() {
		defer writer.Close()
		for _, item := range list {
			if item.Hide {
				continue
			}
			fmt.Fprintf(writer, "%s [%d:%d]\n", item.Name, item.ID, item.Priority)
		}
		utils.RunCMDInSteam(cmdStr, func(line string) {
			index := utils.ArrFindIndex(list, func(item DirItem, _ int) bool {
				return item.Name == line
			})
			if index != -1 {
				return
			}
			item := DirItem{Item{ID: -1, Name: line, Priority: 0, Hide: false}}
			list = append(list, item)
			fmt.Fprintf(writer, "%s [%d:%d]\n", item.Name, item.ID, item.Priority)
		})
	}()

	action, selected, err := utils.RunFZFStream(reader)
	if err != nil {
		if utils.IsCanceled(err) {
			return nil
		}
		return fmt.Errorf("failed to run fzf: %w", err)
	}

	if selected == "" {
		return nil
	}

	index := utils.ArrFindIndex(list, func(item DirItem, _ int) bool {
		return selected == fmt.Sprintf("%s [%d:%d]", item.Name, item.ID, item.Priority)
	})

	if index == -1 {
		return fmt.Errorf("item not found: %s", selected)
	}

	item := list[index]
	if action == utils.Delete {
		item.Hide = true
		if err := dm.Save(item).Error; err != nil {
			return fmt.Errorf("failed to save item: %w", err)
		}
		return JumpDir()
	}
	item.Priority = item.Priority + 1
	if err := dm.Save(item).Error; err != nil {
		return fmt.Errorf("failed to save item: %w", err)
	}

	fmt.Print(`cd `, item.Name)
	return
}

func getDirHistory(dm *dbt.Model) (list []DirItem, err error) {
	oldMap, err := utils.ReadFile("~/.bash_history")
	if err != nil {
		return
	}
	newMap := make(map[string]int)
	for key, v := range oldMap {
		if !strings.HasPrefix(key, "cd") {
			continue
		}
		if strings.Contains(key, "&&") || strings.Contains(key, "../") || strings.Contains(key, "./") {
			continue
		}
		newKey := utils.ExtractPath(key)
		if strings.TrimSpace(newKey) == "" {
			continue
		}
		if !utils.PathExists(newKey) {
			continue
		}
		newMap[newKey] = v

	}
	err = dm.Order("ORDER BY priority DESC").Find(&list).Error
	if err != nil {
		return
	}

	for key, count := range newMap {
		index := utils.ArrFindIndex(list, func(item DirItem, _ int) bool {
			return item.Name == key
		})
		if index != -1 {
			continue
		}
		item := DirItem{Item{-1, key, count, false}}
		list = append(list, item)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Priority > list[j].Priority
	})

	return
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
