package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetCurDirFileName(fileSubFix string) (fileName string, err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	filepath.Base(ex)
	basename := filepath.Base(ex)
	fileName = fmt.Sprintf(`%s.%s`, basename, fileSubFix)
	return
}
func GetCurDirFilePath(fileName string) (filePath string, err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	filepath.Base(ex)
	exPath := filepath.Dir(ex)
	filePath = path.Join(exPath, fileName)
	return
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[1:]), nil
	}
	return path, nil
}

func ReadFile(filePath string) (lineMap map[string]int, err error) {
	filePath, err = expandPath(filePath)
	if err != nil {
		return
	}
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return nil, err
	}
	defer file.Close()

	lineMap = make(map[string]int)
	// 创建一个 Scanner 来逐行读取
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行文本
		lineMap[line]++
	}

	// 检查是否读取中有错误
	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时出错:", err)
		return nil, err
	}

	return
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
