package utils

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

// 专用于捕获fzf输出的函数
func RunFZF(input string) (string, error) {
	cmd := exec.Command("fzf", "--ansi")
	cmd.Stdin = strings.NewReader(input)

	// 创建同时输出到终端和缓冲区的多路写入器
	var buf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf) // 实时显示并捕获
	cmd.Stderr = os.Stderr                       // 错误直接显示

	// 执行命令并等待完成
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// 返回清理后的选择结果
	return strings.TrimSpace(buf.String()), nil
}

// 辅助函数：判断是否用户取消操作
func IsCanceled(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 130
	}
	return false
}
