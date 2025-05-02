package main

import (
	"fmt"
	"go-sqlite-test/dbutils"
	"log"
	"os/exec"
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
	// 打开（或创建）一个数据库文件
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

	rows, err := db.Query("SELECT id, name, priority FROM items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var it Item
		err = rows.Scan(&it.ID, &it.Name, &it.Priority)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, it)
	}

	// 根据 priority 降序排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	println("test:>items.len %d", len(items))

	// 准备 fzf 输入内容
	var fzfInput strings.Builder
	for _, item := range items {
		// 优先级可以显示出来以便参考
		fzfInput.WriteString(fmt.Sprintf("[P%d] %s\n", item.Priority, item.Name))
	}

	// 调用 fzf
	cmd := exec.Command("fzf", "--ansi")
	stdin, _ := cmd.StdinPipe()
	_, _ = cmd.Output()

	// 异步写入 fzf 输入
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(fzfInput.String()))
	}()

	// 获取用户选中结果
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	selected := strings.TrimSpace(string(out))
	fmt.Printf("你选择了: %s\n", selected)
}
