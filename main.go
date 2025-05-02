package main

import (
	"bytes"
	"fmt"
	"go-sqlite-test/dbutils"
	"io"
	"log"
	"os"
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

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// 在一个 goroutine 中读取管道中的数据
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// 执行需要捕获输出的函数
	f()

	// 关闭写入端
	w.Close()
	// 恢复标准输出
	os.Stdout = old
	// 获取捕获的输出
	out := <-outC
	return out
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

	output := captureOutput(func() {

		// 配置并执行fzf命令
		cmd := exec.Command("fzf", "--ansi")

		// 关键修改：绑定标准输入输出到系统终端
		cmd.Stdin = strings.NewReader(fzfInput.String())
		cmd.Stdout = os.Stdout // 直接输出到终端
		cmd.Stderr = os.Stderr

		// 启动命令但不等待完成
		err = cmd.Start()
		if err != nil {
			log.Fatalf("启动fzf失败: %v", err)
		}

		// 等待用户完成选择
		err = cmd.Wait()
	})

	println(":>3----", output)
	// selected := strings.TrimSpace(string(out))
	// if selected == "" {
	// 	log.Println("未选择任何项目")
	// 	return
	// }
	// println(":>4----")

	// fmt.Printf("你选择了: %s\n", selected)
}
