package main

import (
	"fmt"
	"go-sqlite-test/dbt"
	"go-sqlite-test/utils"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := dbt.Init("./example.db")
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
		fzfInput.WriteString(fmt.Sprintf("%s\n", item.Name))
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

	for _, item := range items {
		if item.Name == selected {
			// fmt.Println("你选择了: ", item)
			dbt.UpdateHistoryPriority(db, item.ID, item.Priority+1)
			break
		}

	}
	fmt.Print(selected)

}

// func main() {
// 	// 打开（或创建）数据库
// 	db, err := dbutils.Init("./example.db", "items", `
// 	CREATE TABLE items (
//     	id INTEGER PRIMARY KEY,
//     	name TEXT NOT NULL,
//     	priority INTEGER DEFAULT 0
// 	);`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// 查询数据
// 	items, err := dbutils.GetItems(db)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// 构建fzf输入
// 	var fzfInput strings.Builder
// 	for _, item := range items {
// 		fzfInput.WriteString(fmt.Sprintf("%s\n", item.Name))
// 	}

// 	selected, err := utils.RunFZF(fzfInput.String())
// 	if err != nil {
// 		if utils.IsCanceled(err) { // 检查是否用户取消
// 			log.Println("选择已取消")
// 			return
// 		}
// 		log.Fatal(err)
// 	}

// 	if selected == "" {
// 		log.Println("未选择任何项目")
// 		return
// 	}

// 	for _, item := range items {
// 		if item.Name == selected {
// 			fmt.Println("你选择了: ", item)
// 			dbutils.UpdatePriority(db, item.ID, item.Priority+1)
// 			break
// 		}

// 	}
// }
