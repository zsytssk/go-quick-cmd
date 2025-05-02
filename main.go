package main

import (
	"fmt"
	"go-sqlite-test/dbutils"
	"go-sqlite-test/utils"
	"log"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID       int
	Name     string
	Priority int
}

func main() {
	// 打开（或创建）数据库
	db, err := dbutils.Init("./example.db", "items", `
	CREATE TABLE items (
    	id INTEGER PRIMARY KEY,
    	name TEXT NOT NULL,
    	priority INTEGER DEFAULT 0
	);`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 查询数据
	rows, err := db.Query("SELECT id, name, priority FROM items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 读取数据到结构体切片
	var items []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Name, &it.Priority); err != nil {
			log.Fatal(err)
		}
		items = append(items, it)
	}

	// 按优先级降序排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	// 构建fzf输入
	var fzfInput strings.Builder
	for _, item := range items {
		fzfInput.WriteString(fmt.Sprintf("[P%d] %s\n", item.Priority, item.Name))
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

	fmt.Printf("你选择了: %s\n", selected)

}
