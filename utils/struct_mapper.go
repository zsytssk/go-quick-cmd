package utils

import "reflect"

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)

	// 如果是指针，先取指针指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		fieldVal := v.Field(i)
		if !fieldVal.CanInterface() {
			continue
		}
		value := fieldVal.Interface()
		if IsZero(value) {
			continue
		}
		if fieldType.Type.Kind() == reflect.Struct {
			result[fieldType.Name] = StructToMap(value)
			continue
		}
		result[fieldType.Name] = value
	}

	return result
}

func MapToStruct(obj map[string]interface{}, t reflect.Type) interface{} {
	// 创建结构体实例
	v := reflect.New(t).Elem()

	for key, value := range obj {
		// 获取结构体字段
		field := v.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		if field.Type().Kind() == reflect.Struct {
			// 递归处理嵌套结构体
			nestedMap, ok := value.(map[string]interface{})
			if !ok {
				continue
			}
			nestedStruct := MapToStruct(nestedMap, field.Type())

			field.Set(reflect.ValueOf(nestedStruct))
			continue
		}
		// 设置字段值
		field.Set(reflect.ValueOf(value))
	}

	return v.Interface()
}
