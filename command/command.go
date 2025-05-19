package command

import (
	"database/sql"
	"fmt"
	"quick-cmd/dbt"
	"quick-cmd/utils"
)

// Command 定义了命令接口
type Command interface {
	Execute() error
	GetName() string
}

// BaseCommand 提供了基础命令实现
type BaseCommand struct {
	Name string
	DB   *sql.DB
}

// GetName 返回命令名称
func (b *BaseCommand) GetName() string {
	return b.Name
}

// NewBaseCommand 创建基础命令
func NewBaseCommand(name string) (*BaseCommand, error) {
	dbPath, err := utils.GetCurDirFileName("db")
	if err != nil {
		return nil, fmt.Errorf("failed to get db path: %w", err)
	}

	db, err := dbt.Init(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return &BaseCommand{
		Name: name,
		DB:   db,
	}, nil
}
