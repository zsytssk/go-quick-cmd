package dbutils

import (
	"database/sql"
	"fmt"
	"log"
)

func Init(dbPath string, tableName string, tableFmt string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}
	exist := checkTableExist(db, tableName)
	if !exist {
		createTable(db, tableFmt)
		insertTestData(db)
	}
	return
}

func checkTableExist(db *sql.DB, tableName string) (exist bool) {
	// 执行查询，检查 sqlite_master 表中是否存在指定名称的表
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
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
