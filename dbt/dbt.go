package dbt

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"quick-cmd/utils"
	"reflect"
	"strings"
)

type TableStruct interface {
	TableName() string
}

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

func StructToSQLCreateTable(db *sql.DB, obj TableStruct, fields_list []FieldItem) (err error) {
	if CheckTableExist(db, obj.TableName()) {
		return
	}
	var columns []string
	for _, field := range fields_list {
		db_type := field.DbType
		if db_type == "primaryKey" {
			columns = append(columns, fmt.Sprintf("  %s %s PRIMARY KEY", field.Name, field.SqlType))
			continue
		}
		columns = append(
			columns,
			fmt.Sprintf("  %s %s NOT NULL DEFAULT %s",
				field.Name,
				field.SqlType,
				formatSQLDefaultValue(field.OriType),
			))
	}

	sqlStr := fmt.Sprintf("CREATE TABLE %s (\n%s\n);",
		strings.ToLower(obj.TableName()),
		strings.Join(columns, ",\n"),
	)
	// fmt.Println(sqlStr)
	_, err = db.Exec(sqlStr)
	if err != nil {
		return
	}
	return
}
func StructToSQLInsert(db *sql.DB, obj TableStruct, fields_list []FieldItem) (err error) {
	exists, err := CheckItemExist(db, obj, fields_list)
	if exists || err != nil {
		return
	}
	var columns []string
	var placeholders []string
	var values []interface{}
	for _, field := range fields_list {
		db_type := field.DbType
		if db_type == "primaryKey" {
			continue
		}
		columns = append(columns, field.Name)
		placeholders = append(placeholders, "?")
		values = append(values, field.Value)
	}
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		strings.ToLower(obj.TableName()),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	_, err = db.Exec(sqlStr, values...)
	if err != nil {
		return
	}
	return
}
func StructToSQLUpdate(
	db *sql.DB,
	obj TableStruct,
	fields_list []FieldItem,
	ignore_zero bool,
) (err error) {
	var columns []string
	var where_str string
	for _, field := range fields_list {
		db_type := field.DbType
		if db_type == "primaryKey" {
			where_str = fmt.Sprintf(" %s = %v", field.Name, formatSQLValue(field.Value))
			continue
		}
		if ignore_zero && utils.IsZero(field.Value) {
			continue
		}
		columns = append(columns, fmt.Sprintf("  %s = %v", field.Name, formatSQLValue(field.Value)))
	}
	sqlStr := fmt.Sprintf("UPDATE %s\nSET %s \n WHERE %s;",
		strings.ToLower(obj.TableName()),
		strings.Join(columns, ",\n"),
		where_str,
	)
	// fmt.Println(sqlStr)
	_, err = db.Exec(sqlStr)
	if err != nil {
		return
	}
	return
}

func StructToSQLDelete(db *sql.DB, obj TableStruct, fields_list []FieldItem) (err error) {
	// DELETE FROM users WHERE id = 3;
	var where_str string
	for _, field := range fields_list {
		db_type := field.DbType
		if db_type == "primaryKey" {
			where_str = fmt.Sprintf("  %s = %v", field.Name, formatSQLValue(field.Value))
			break
		}
	}
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s;",
		strings.ToLower(obj.TableName()),
		where_str,
	)
	_, err = db.Exec(sqlStr)
	if err != nil {
		return
	}
	return
}

func StructToSQLGetList(db *sql.DB, obj TableStruct) (list []TableStruct, err error) {
	fields_list, err := collectFields(obj)
	if err != nil {
		return
	}
	var columns []string
	for _, field := range fields_list {
		columns = append(columns, field.Name)
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(columns, ", "),
		strings.ToLower(obj.TableName()),
	)
	rows, err := db.Query(sqlStr)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		elemPtr := reflect.New(reflect.TypeOf(obj)) // *T
		elemVal := elemPtr.Elem()                   // T
		var scanArgs []interface{}

		for _, field := range fields_list {
			fieldName := field.OriName

			// 获取结构体中对应的字段
			structField := elemVal.FieldByName(fieldName) // 需要转换

			// 跳过非法或未导出字段
			if !structField.IsValid() || !structField.CanAddr() {
				continue
			}

			// 添加字段地址作为 Scan 参数
			scanArgs = append(scanArgs, structField.Addr().Interface())
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, elemVal.Interface().(TableStruct))
	}
	return
}

func CheckItemExist(db *sql.DB, obj TableStruct, fields_list []FieldItem) (exists bool, err error) {
	var where_str string
	for _, field := range fields_list {
		db_type := field.DbType
		if db_type == "primaryKey" {
			where_str = fmt.Sprintf("  %s = %v", field.Name, formatSQLValue(field.Value))
			break
		}
	}

	sqlStr := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s);",
		strings.ToLower(obj.TableName()),
		where_str,
	)

	err = db.QueryRow(sqlStr).Scan(&exists)
	return
}

func StructToSQLDeleteTable(db *sql.DB, obj TableStruct) (err error) {
	// DROP TABLE IF EXISTS table_name;
	sqlStr := fmt.Sprintf("DROP TABLE IF EXISTS %s;",
		strings.ToLower(obj.TableName()),
	)
	_, err = db.Exec(sqlStr)
	if err != nil {
		return
	}
	return
}

type FieldItem struct {
	DbType  string
	Name    string
	OriName string
	SqlType string
	OriType reflect.Type
	Value   interface{}
}

func collectFields(obj interface{}) (connects []FieldItem, err error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		errors.New("这是一个简单的错误")
		return
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldVal := v.Field(i)
		// 跳过未导出字段
		if !fieldVal.CanInterface() {
			continue
		}
		value := fieldVal.Interface()
		// 匿名字段 && 是 struct，递归处理
		if fieldType.Anonymous && fieldType.Type.Kind() == reflect.Struct {
			sub_list, err := collectFields(value)
			if err != nil {
				return connects, err
			}
			connects = append(connects, sub_list...)
			continue
		}

		// 字段名使用 json tag，否则用字段名
		fieldName := fieldType.Tag.Get("json")
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = utils.CamelToSnake(fieldType.Name)
		}
		dbType := fieldType.Tag.Get("db")

		sqlType := GoTypeToSQLType(fieldType.Type, dbType == "primaryKey")
		connects = append(connects, FieldItem{
			DbType:  dbType,
			Name:    fieldName,
			OriName: fieldType.Name,
			OriType: fieldType.Type,
			SqlType: sqlType,
			Value:   value,
		})
	}

	return
}

type TableColumn struct {
	Cid       int
	Name      string
	Ctype     string
	Notnull   int
	DfltValue sql.NullString
	Pk        int
}

func getTableColumns(db *sql.DB, obj TableStruct) (columns []TableColumn, err error) {
	sqlStr := fmt.Sprintf("PRAGMA table_info(%s);", strings.ToLower(obj.TableName()))
	rows, err := db.Query(sqlStr)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var column TableColumn

		err = rows.Scan(
			&column.Cid,
			&column.Name,
			&column.Ctype,
			&column.Notnull,
			&column.DfltValue,
			&column.Pk,
		)
		if err != nil {
			return
		}
		columns = append(columns, column)
	}
	return
}

func CheckTableExist(db *sql.DB, tableName string) (exist bool) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name=?)"
	err := db.QueryRow(query, tableName).Scan(&exists)
	return err == nil && exists
}

func SyncTableColumns(db *sql.DB, obj TableStruct, fields_list []FieldItem) (err error) {
	columns, err := getTableColumns(db, obj)
	if err != nil {
		return
	}

	for _, column := range columns {
		if utils.ArrFindIndex(fields_list, func(field FieldItem, index int) bool {
			return field.Name == column.Name &&
				IsSQLTypeCompatible(column.Ctype, field.OriType)
		}) != -1 {
			continue
		}
		// fmt.Println("deleteColumns:>2", column.Name, column.Ctype)
		_, err = db.Exec(fmt.Sprintf(
			"ALTER TABLE %s DROP COLUMN %s;",
			obj.TableName(),
			column.Name,
		))
		if err != nil {
			return
		}
	}
	for _, field := range fields_list {
		if utils.ArrFindIndex(columns, func(column TableColumn, index int) bool {
			return column.Name == field.Name &&
				IsSQLTypeCompatible(column.Ctype, field.OriType)
		}) != -1 {
			continue
		}
		// fmt.Println("deleteColumns:>2", field.Name, field.SqlType)

		_, err = db.Exec(fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN %s %s NOT NULL DEFAULT %s;",
			obj.TableName(),
			field.Name,
			field.SqlType,
			formatSQLDefaultValue(field.OriType),
		))
	}

	return
}
