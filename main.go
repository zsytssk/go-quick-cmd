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

	"github.com/creack/pty" // 新增pty支持
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

	// 配置并执行fzf命令
	cmd := exec.Command("fzf", "--ansi")
	cmd.Stdin = strings.NewReader(fzfInput.String())

	// 使用伪终端启动进程
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatalf("启动失败: %v", err)
	}
	defer ptmx.Close()

	// 实时捕获输出
	var output bytes.Buffer
	go func() {
		io.Copy(io.MultiWriter(&output, os.Stdout), ptmx) // 同时输出到终端和缓冲区
	}()

	// 等待命令完成
	err = cmd.Wait()

	// 错误处理
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			switch exitErr.ExitCode() {
			case 130:
				log.Println("选择已取消")
				return
			case 1:
				log.Println("没有选择任何项目")
				return
			}
		}
		log.Fatalf("执行错误: %v", err)
	}

	// 处理输出
	selected := strings.TrimSpace(output.String())
	fmt.Printf("\n你选择了: %s\n", selected)
}
