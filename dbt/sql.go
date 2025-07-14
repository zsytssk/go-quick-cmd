package dbt

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func formatSQLValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		// 空字符串也必须写成 ''
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case time.Time:
		// 注意：Time 应该加引号包住
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case nil:
		return "NULL"
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		// 数值类型直接返回字符串形式（不加引号）
		return fmt.Sprintf("%v", v)
	}
}

func formatSQLDefaultValue(t reflect.Type) string {
	k := t.Kind()
	switch {
	case k == reflect.String:
		return "''"
	case isIntKind(k), isUintKind(k), isFloatKind(k):
		return "0"
	case k == reflect.Bool:
		return "false"
	default:
		return "NULL"
	}
}

func GoTypeToSQLType(goType reflect.Type, isPrimary bool) string {
	switch goType.Kind() {
	case reflect.Int, reflect.Int32:
		if isPrimary {
			return "INTEGER"
		}
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint64:
		return "BIGINT UNSIGNED"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.Bool:
		return "TINYINT(1)"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Slice:
		if goType.Elem().Kind() == reflect.Uint8 {
			return "BLOB"
		}
	case reflect.Struct:
		if goType.PkgPath() == "time" && goType.Name() == "Time" {
			return "DATETIME"
		}
	}
	return "TEXT"
}

func IsSQLTypeCompatible(sqlType string, goType reflect.Type) bool {
	sqlType = strings.ToUpper(sqlType)

	switch sqlType {
	case "INTEGER", "INT", "BIGINT", "TINYINT":
		return isIntKind(goType.Kind())
	case "FLOAT", "DOUBLE", "REAL", "DECIMAL":
		return isFloatKind(goType.Kind())
	case "BOOLEAN", "BOOL", "TINYINT(1)":
		return goType.Kind() == reflect.Bool
	case "VARCHAR", "CHAR", "TEXT", "LONGTEXT", "MEDIUMTEXT":
		return goType.Kind() == reflect.String
	case "DATE", "DATETIME", "TIMESTAMP":
		return goType.PkgPath() == "time" && goType.Name() == "Time"
	case "BLOB", "BYTEA":
		return goType.Kind() == reflect.Slice && goType.Elem().Kind() == reflect.Uint8
	default:
		// 默认用 string 类型兜底
		return goType.Kind() == reflect.String
	}
}

func isIntKind(k reflect.Kind) bool {
	return k >= reflect.Int && k <= reflect.Int64
}

func isUintKind(k reflect.Kind) bool {
	return k >= reflect.Uint && k <= reflect.Uintptr
}

func isFloatKind(k reflect.Kind) bool {
	return k == reflect.Float32 || k == reflect.Float64
}
