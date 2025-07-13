package command

import (
	"database/sql"
	"fmt"
	"quick-cmd/dbt"
	"quick-cmd/utils"
)

type Item struct {
	ID       int    `db:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Hide     bool   `json:"hidden"`
}

func getDb() (db *sql.DB, err error) {
	dbPath, err := utils.GetCurDirFileName("db")
	if err != nil {
		return nil, fmt.Errorf("failed to get db path: %w", err)
	}

	db, err = dbt.Init(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}
	return
}
