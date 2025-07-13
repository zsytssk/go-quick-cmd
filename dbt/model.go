package dbt

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

type Condition struct {
	ConditionType string
	Query         string
	Args          []interface{}
}
type TableInfo struct {
	FieldsList []map[string]interface{}
	TableName  string
}

type Model struct {
	DB         *sql.DB
	Obj        TableStruct
	Error      error
	TableInfo  TableInfo
	conditions []Condition
}

func NewModel(db *sql.DB, obj TableStruct) *Model {
	fields_list := collectFields(obj)
	m := &Model{
		DB:        db,
		Obj:       obj,
		TableInfo: TableInfo{FieldsList: fields_list, TableName: obj.TableName()},
	}
	if !CheckTableExist(m.DB, m.Obj.TableName()) {
		err := StructToSQLCreateTable(m.DB, m.Obj)
		if err != nil {
			m.Error = err
			return m
		}
	}
	return m
}

func (m *Model) BuildConditions(conditions []Condition) string {
	var where_str string
	var order_str string
	var limit_str string
	var offset_str string
	for _, condition := range conditions {
		switch condition.ConditionType {
		case "where":
			if where_str != "" {
				continue
			}
			where_str = condition.Query
		case "order":
			if order_str != "" {
				continue
			}
			order_str = condition.Query
		case "limit":
			if limit_str != "" {
				continue
			}
			limit_str = condition.Query
		case "offset":
			if offset_str != "" {
				continue
			}
			offset_str = condition.Query
		}
	}
	return fmt.Sprintf("%s %s %s %s %s",
		m.TableInfo.TableName,
		where_str,
		order_str,
		limit_str,
		offset_str,
	)
}

func (m *Model) Where(query string) *Model {
	m.conditions = append(m.conditions, Condition{
		ConditionType: "where",
		Query:         query,
	})
	return m
}

func (m *Model) Limit(limit int64) *Model {
	m.conditions = append(m.conditions, Condition{
		ConditionType: "limit",
		Query:         fmt.Sprintf(" LIMIT %d", limit),
	})
	return m
}
func (m *Model) Offset(offset int64) *Model {
	m.conditions = append(m.conditions, Condition{
		ConditionType: "offset",
		Query:         fmt.Sprintf(" OFFSET %d", offset),
	})
	return m
}
func (m *Model) Order(order string) *Model {
	m.conditions = append(m.conditions, Condition{
		ConditionType: "order",
		Query:         order,
	})
	return m
}

func (m *Model) Save(value interface{}) *Model {
	exists, err := CheckItemExist(m.DB, value.(TableStruct))
	if err != nil {
		m.Error = err
		return m
	}
	if !exists {
		err = StructToSQLInsert(m.DB, value.(TableStruct))
		m.Error = err
		return m
	}
	err = StructToSQLUpdate(m.DB, value.(TableStruct), false)
	m.Error = err
	return m
}

func (m *Model) Update(value interface{}) *Model {
	err := StructToSQLUpdate(m.DB, value.(TableStruct), true)
	m.Error = err
	return m
}

func (m *Model) Delete(value interface{}) *Model {
	err := StructToSQLDelete(m.DB, value.(TableStruct))
	m.Error = err
	return m
}

func (m *Model) First(first interface{}) *Model {
	err := SyncTableColumns(m.DB, m.Obj)
	if err != nil {
		m.Error = err
		return m
	}
	conditions := append([]Condition{{
		ConditionType: "order",
		Query:         "LIMIT 1",
	}}, m.conditions...)

	sql := fmt.Sprintf("SELECT * FROM %s", m.BuildConditions(conditions))
	elemPtr := reflect.ValueOf(first) // *T
	elemVal := elemPtr.Elem()         // T
	var scanArgs []interface{}

	for _, field := range m.TableInfo.FieldsList {
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
	err = m.DB.QueryRow(sql).Scan(scanArgs...)
	m.Error = err
	return m
}
func (m *Model) Count(count *int64) *Model {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.BuildConditions(m.conditions))
	err := m.DB.QueryRow(
		sql,
	).Scan(count)
	m.Error = err
	return m
}
func (m *Model) Find(dest interface{}) *Model {
	err := SyncTableColumns(m.DB, m.Obj)
	if err != nil {
		m.Error = err
		return m
	}
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() == reflect.Ptr {
		destVal = destVal.Elem() // 解引用 *slice => slice
	}

	sql := fmt.Sprintf("SELECT * FROM %s", m.BuildConditions(m.conditions))
	// fmt.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		m.Error = err
		return m
	}
	defer rows.Close()
	for rows.Next() {
		typ := reflect.TypeOf(m.Obj)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem() // 获取指针指向的真实类型
		}
		elemPtr := reflect.New(typ) // *T
		elemVal := elemPtr.Elem()   // T
		var scanArgs []interface{}

		for _, field := range m.TableInfo.FieldsList {
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
		destVal = reflect.Append(destVal, elemVal)
	}
	reflect.ValueOf(dest).Elem().Set(destVal)

	return m
}
