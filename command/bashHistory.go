package command

import (
	"fmt"
	"log"
	"quick-cmd/dbt"
	"quick-cmd/utils"
	"strings"
)

func HashHistory() {
	dbPath, err := utils.GetCurDirFileName("db")
	if err != nil {
		log.Fatal(err)
	}
	db, err := dbt.Init(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	items, err := dbt.GetHistory(db)
	if err != nil {
		log.Fatal(err)
	}

	// 构建fzf输入
	var fzfInput strings.Builder
	for _, item := range items {
		fzfInput.WriteString(fmt.Sprintf("%s [%d:%d]\n", item.Name, item.ID, item.Priority))
	}

	selected, err := utils.RunFZF(fzfInput.String())
	if err != nil {
		if utils.IsCanceled(err) { // 检查是否用户取消
			log.Println("选择已取消")
			return
		}
		log.Fatal(err)
	}

	if selected == "" {
		log.Println("未选择任何项目")
		return
	}

	index := utils.ArrFindIndex(items, func(item dbt.Item, _ int) bool {
		return selected == fmt.Sprintf("%s [%d:%d]", item.Name, item.ID, item.Priority)
	})

	if index == -1 {
		log.Println("find item index = ", selected)
		return
	}
	item := items[index]
	dbt.UpdateHistoryPriority(db, item)
	fmt.Print(item.Name)
}
