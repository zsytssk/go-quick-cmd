package main

import (
	"database/sql"
	"fmt"
	"go-sqlite-test/dbutils"
	"go-sqlite-test/utils"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}

	lineMap, err := utils.ReadFile("~/.bash_history")
	if err != nil {
		log.Fatal(err)
	}
	dbutils.InitHistory(db, lineMap)

	// 	// 查询数据
	items, err := dbutils.GetHistory(db)
	if err != nil {
		log.Fatal(err)
	}

	// 构建fzf输入
	var fzfInput strings.Builder
	for _, item := range items {
		fzfInput.WriteString(fmt.Sprintf("%s\n", item.Name))
	}

	fmt.Print("--->1")
	selected, err := utils.RunFZF(fzfInput.String())
	fmt.Print("--->1")
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

	fmt.Print("--->2")
	for _, item := range items {
		if item.Name == selected {
			// fmt.Println("你选择了: ", item)
			dbutils.UpdateHistoryPriority(db, item.ID, item.Priority+1)
			break
		}

	}
	// fmt.Print("--->3", selected)

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
