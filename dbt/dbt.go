package dbt

import (
	"database/sql"
	"fmt"
	"log"
	"quick-cmd/utils"
	"reflect"
	"strings"
)

type TableInterface interface {
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

type TableStruct interface {
	TableName() string
}

func StructToSQLCreateTable(db *sql.DB, obj TableStruct) (err error) {
	if CheckTableExist(db, obj.TableName()) {
		return
	}
	fields_list := collectFields(obj)
	var columns []string
	for _, field := range fields_list {
		db_type := field["dbType"].(string)
		if db_type == "primaryKey" {
			columns = append(columns, fmt.Sprintf("  %s %s PRIMARY KEY", field["name"], field["sqlType"]))
			continue
		}
		columns = append(columns, fmt.Sprintf("  %s %s", field["name"], field["sqlType"]))
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
func StructToSQLInsert(db *sql.DB, obj TableStruct) (err error) {
	exists, err := CheckItemExist(db, obj)
	if exists || err != nil {
		return
	}
	fields_list := collectFields(obj)
	var columns []string
	var placeholders []string
	var values []interface{}
	for _, field := range fields_list {
		db_type := field["dbType"].(string)
		if db_type == "primaryKey" {
			continue
		}
		columns = append(columns, field["name"].(string))
		placeholders = append(placeholders, "?")
		values = append(values, field["value"])
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
func StructToSQLUpdate(db *sql.DB, obj TableStruct, ignore_zero bool) (err error) {
	fields_list := collectFields(obj)
	var columns []string
	var where_str string
	for _, field := range fields_list {
		db_type := field["dbType"].(string)
		if db_type == "primaryKey" {
			where_str = fmt.Sprintf(" %s = %v", field["name"], formatSQLValue(field["value"]))
			continue
		}
		if ignore_zero && utils.IsZero(field["value"]) {
			continue
		}
		columns = append(columns, fmt.Sprintf("  %s = %v", field["name"], formatSQLValue(field["value"])))
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

func StructToSQLDelete(db *sql.DB, obj TableStruct) (err error) {
	// DELETE FROM users WHERE id = 3;
	fields_list := collectFields(obj)
	var where_str string
	for _, field := range fields_list {
		db_type := field["dbType"].(string)
		if db_type == "primaryKey" {
			where_str = fmt.Sprintf("  %s = %v", field["name"], formatSQLValue(field["value"]))
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

func StructToSQLGetList[T TableStruct](db *sql.DB, obj T) (list []T, err error) {
	fields_list := collectFields(obj)
	var columns []string
	for _, field := range fields_list {
		columns = append(columns, fmt.Sprintf("%s", field["name"]))
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
			fieldName := field["oriName"].(string)

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
		list = append(list, elemVal.Interface().(T))
	}
	return
}

func CheckItemExist(db *sql.DB, obj TableStruct) (exists bool, err error) {
	fields_list := collectFields(obj)
	var where_str string
	for _, field := range fields_list {
		db_type := field["dbType"].(string)
		if db_type == "primaryKey" {
			where_str = fmt.Sprintf("  %s = %v", field["name"], formatSQLValue(field["value"]))
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
	Value   interface{}
}

func collectFields(obj interface{}) (connects []map[string]interface{}) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldVal := v.Field(i)

		// 跳过未导出字段
		if !fieldType.IsExported() || !fieldVal.CanInterface() {
			continue
		}
		value := fieldVal.Interface()
		// 匿名字段 && 是 struct，递归处理
		if fieldType.Anonymous && fieldType.Type.Kind() == reflect.Struct {
			sub_list := collectFields(value)
			connects = append(connects, sub_list...)
			continue
		}

		// 字段名使用 json tag，否则用字段名
		name := fieldType.Tag.Get("json")
		if name == "" || name == "-" {
			name = utils.CamelToSnake(fieldType.Name)
		}
		dbType := fieldType.Tag.Get("db")

		sqlType := GoTypeToSQLType(fieldType.Type)
		connects = append(connects, map[string]interface{}{
			"dbType":  dbType,
			"name":    name,
			"oriName": fieldType.Name,
			"oriType": fieldType.Type,
			"sqlType": sqlType,
			"value":   value,
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
		log.Fatal(err)
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

func SyncTableColumns(db *sql.DB, obj TableStruct) (err error) {
	columns, err := getTableColumns(db, obj)
	fields_list := collectFields(obj)

	for _, column := range columns {
		if utils.ArrFindIndex(fields_list, func(field map[string]interface{}, index int) bool {
			return field["name"] == column.Name &&
				IsSQLTypeCompatible(column.Ctype, field["oriType"].(reflect.Type))
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
			return column.Name == field["name"] &&
				IsSQLTypeCompatible(column.Ctype, field["oriType"].(reflect.Type))
		}) != -1 {
			continue
		}
		// fmt.Println("deleteColumns:>2", field["name"], field["sqlType"])

		_, err = db.Exec(fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN %s %s NOT NULL DEFAULT %s;",
			obj.TableName(),
			field["name"],
			field["sqlType"],
			formatSQLDefaultValue(field["oriType"].(reflect.Type)),
		))
	}

	return
}
