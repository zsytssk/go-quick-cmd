package dbutils

import (
	"database/sql"
	"fmt"
	"log"
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
func UpdateItemPriority(db *sql.DB, tableName string, id int, priority int) error {
	stmt := fmt.Sprintf(`UPDATE %s
         SET priority = ?
         WHERE id = ?
         AND priority <> ?`, tableName)

	result, err := db.Exec(
		stmt,
		priority, id, priority,
	)
	if err != nil {
		return fmt.Errorf("更新失败: %w", err)
	}

	// 简化版影响检查
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("记录不存在或值未变化")
	}

	return nil
}

type Item struct {
	ID       int
	Name     string
	Priority int
}

func GetItems(db *sql.DB, tableName string) (items []Item, err error) {
	queryStr := fmt.Sprintf("SELECT id, name, priority FROM %s ORDER BY priority desc", tableName)
	rows, err := db.Query(queryStr)
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

	return
}
