package utils

import (
	"os"
	"regexp"
	"strings"
)

func GetCmd() *string {
	if len(os.Args) == 1 {
		return nil
	}
	return &os.Args[1]
}

func ExtractPath(input string) string {
	// 1. 移除 "cd" 及其后的所有空格
	re := regexp.MustCompile(`^cd\s*`)
	path := re.ReplaceAllString(input, "")

	// 2. 去除首尾空白字符
	path = strings.TrimSpace(path)

	// 3. 检查并去除包裹路径的单引号或双引号
	if len(path) >= 2 {
		quoteChar := path[0]
		if (quoteChar == '\'' || quoteChar == '"') && path[len(path)-1] == quoteChar {
			path = path[1 : len(path)-1]
		}
	}

	return path
}
