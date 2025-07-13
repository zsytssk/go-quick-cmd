package utils

import (
	"regexp"
	"strings"
)

func CamelToSnake(s string) string {
	// 替换所有大写字母前加下划线
	snake := regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
