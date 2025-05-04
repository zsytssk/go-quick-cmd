package dbt

import (
	"database/sql"
	"fmt"
	"go-sqlite-test/utils"
	"sort"
	"strings"
)

func InitDirTable(db *sql.DB, lineMap map[string]int) {
	exist := checkTableExist(db, "dir")
	if !exist {
		createTable(db, genTableStmt("dir"))
	}

	stmtSelect := `SELECT priority FROM dir WHERE name = ?`
	stmtInsert := `INSERT INTO dir (name, priority) VALUES (?, ?)`
	stmtUpdate := `UPDATE dir SET priority = ? WHERE name = ?`

	for k, v := range lineMap {

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

func GetDir(db *sql.DB) (items []Item, err error) {
	oldMap, err := utils.ReadFile("~/.bash_history")

	if err != nil {
		return
	}
	newMap := make(map[string]int)
	for key, v := range oldMap {
		if !strings.HasPrefix(key, "cd") {
			continue
		}
		if strings.Contains(key, "&&") || strings.Contains(key, "../") || strings.Contains(key, "./") {
			continue
		}
		newKey := utils.ExtractPath(key)
		if strings.TrimSpace(newKey) == "" {
			continue
		}
		newMap[newKey] = v

	}

	exist := checkTableExist(db, "dir")
	if !exist {
		InitDirTable(db, newMap)
	}

	items, err = GetItems(db, "dir")
	if err != nil {
		return
	}
	for key, count := range newMap {
		index := utils.ArrFindIndex(items, func(item Item, _ int) bool {
			return item.Name == key
		})
		if index != -1 {
			continue
		}
		item := Item{-1, key, count}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	return
}

func UpdateDirPriority(db *sql.DB, item Item) (err error) {
	if item.ID != -1 {
		return UpdateItemPriority(db, "dir", item.ID, item.Priority+1)
	}
	return InsertItemPriority(db, "dir", item.Name, item.Priority+1)
}
