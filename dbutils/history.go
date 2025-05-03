package dbutils

import (
	"database/sql"
	"fmt"
)

func InitHistory(db *sql.DB, lineMap map[string]int) {
	exist := checkTableExist(db, "history")
	if !exist {
		createTable(db, genTableStmt("history"))
	}

	stmtSelect := `SELECT priority FROM history WHERE name = ?`
	stmtInsert := `INSERT INTO history (name, priority) VALUES (?, ?)`
	stmtUpdate := `UPDATE history SET priority = ? WHERE name = ?`

	for k, v := range lineMap {
		if v <= 2 {
			continue // 只处理出现次数大于 2 的
		}

		var currentPriority int
		err := db.QueryRow(stmtSelect, k).Scan(&currentPriority)

		switch {
		case err == sql.ErrNoRows:
			// 不存在，插入新记录
			_, _ = db.Exec(stmtInsert, k, v)

		case err == nil && currentPriority < v:
			// 已存在，且 priority 更小，更新
			_, _ = db.Exec(stmtUpdate, v, k)

		// err == nil 且 currentPriority >= v，不处理
		case err != nil:
			// 其他查询错误
			fmt.Println("查询失败:", err)
		}
	}
}

func GetHistory(db *sql.DB) (items []Item, err error) {
	return GetItems(db, "history")
}
func UpdateHistoryPriority(db *sql.DB, id int, priority int) (err error) {
	return UpdateItemPriority(db, "history", id, priority)
}
