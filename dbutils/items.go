package dbutils

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
)

func insertTestData(db *sql.DB) {
	items := []struct {
		name     string
		priority int
	}{
		{"Learn Go", 5},
		{"Do laundry", 2},
		{"Buy milk", 1},
		{"Fix bug", 10},
	}

	stmt, err := db.Prepare("INSERT INTO items(name, priority) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, item := range items {
		_, err = stmt.Exec(item.name, item.priority)
		if err != nil {
			log.Printf("插入失败: %v\n", err)
		}
	}
}

// UpdatePriority 根据名称更新优先级
func UpdatePriority(db *sql.DB, name string, newPriority int) error {
	// 参数校验
	if name == "" {
		return fmt.Errorf("名称不能为空")
	}

	// 使用事务保证原子性
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // 安全回滚

	// 执行更新
	result, err := tx.Exec(
		"UPDATE items SET priority = ? WHERE id = ?",
		newPriority,
		name,
	)
	if err != nil {
		return fmt.Errorf("更新失败: %w", err)
	}

	// 检查影响行数
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("名称 '%s' 不存在", name)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

type Item struct {
	ID       int
	Name     string
	Priority int
}

func GetItems(db *sql.DB) (items []Item, err error) {
	rows, err := db.Query("SELECT id, name, priority FROM items")
	if err != nil {
		return
	}
	defer rows.Close()

	// 读取数据到结构体切片
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
	return
}
