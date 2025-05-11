package dbt

import (
	"log"
	"quick-cmd/utils"

	"database/sql"
	"fmt"
)

func Init(dbPath string) (db *sql.DB, err error) {
	filePath, err := utils.GetCurDirFilePath(dbPath)
	if err != nil {
		return
	}
	db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return
	}
	return
}

type Item struct {
	ID       int
	Name     string
	Priority int
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

// UpdatePriority 根据名称更新优先级
func InsertItemPriority(db *sql.DB, tableName string, name string, priority int) error {
	stmt := fmt.Sprintf(`INSERT INTO %s (name, priority) VALUES (?, ?)`, tableName)

	result, err := db.Exec(
		stmt,
		name, priority,
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

func checkTableExist(db *sql.DB, tableName string) (exist bool) {
	// 执行查询，检查 sqlite_master 表中是否存在指定名称的表
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return false
	}

	// 如果查询结果大于 0，则表示表存在
	return count > 0
}

func createTable(db *sql.DB, sqlStmt string) (sus bool, err error) {
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return false, err
	}
	return true, nil
}

func genTableStmt(tableName string) string {
	return fmt.Sprintf(`
	CREATE TABLE %s (
    	id INTEGER PRIMARY KEY,
    	name TEXT NOT NULL,
    	priority INTEGER DEFAULT 0
	);`, tableName)
}
