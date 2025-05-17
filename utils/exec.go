package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

func RunFZF(input string) (string, error) {
	// 创建伪终端
	ptm, pts, err := pty.Open()
	if err != nil {
		return "", err
	}
	defer ptm.Close()
	defer pts.Close()

	// 创建结果缓冲区
	// 配置fzf命令
	var buf bytes.Buffer
	fzf := exec.Command("fzf", "--ansi")
	fzf.Stdout = io.MultiWriter(pts, &buf) // 实时显示并捕获

	fzf.Stderr = os.Stderr
	fzf.Stdin = strings.NewReader(input) // 允许接收键盘输入

	// 执行命令并等待完成
	if err := fzf.Run(); err != nil {
		return "", err
	}

	// 返回清理后的结果
	return strings.TrimSpace(buf.String()), nil
}

func RunFZFStream(reader io.Reader) (string, error) {
	// 创建伪终端
	ptm, pts, err := pty.Open()
	if err != nil {
		return "", err
	}
	defer ptm.Close()
	defer pts.Close()

	// 创建结果缓冲区
	// 配置fzf命令
	var buf bytes.Buffer
	fzf := exec.Command("fzf", "--ansi")
	fzf.Stdout = io.MultiWriter(pts, &buf) // 实时显示并捕获

	fzf.Stderr = os.Stderr
	fzf.Stdin = reader // 允许接收键盘输入

	// 执行命令并等待完成
	if err := fzf.Run(); err != nil {
		return "", err
	}

	// 返回清理后的结果
	return strings.TrimSpace(buf.String()), nil
}

func RunCMD(input string) (string, error) {
	// 创建伪终端
	ptm, pts, err := pty.Open()
	if err != nil {
		return "", err
	}
	defer ptm.Close()
	defer pts.Close()

	// 创建结果缓冲区
	// 配置fzf命令
	var buf bytes.Buffer
	cmd := exec.Command("bash", "-c", input)
	cmd.Stdout = io.MultiWriter(pts, &buf) // 实时显示并捕获

	cmd.Stderr = os.Stderr
	// cmd.Stdin = io.MultiReader(strings.NewReader(input)) // 允许接收键盘输入

	// 执行命令并等待完成
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// 返回清理后的结果
	return strings.TrimSpace(buf.String()), nil
}

func RunCMDInSteam(input string, fn func(string)) {
	cmd := exec.Command("bash", "-c", input)
	cmd.Stderr = os.Stderr
	cmdOut, _ := cmd.StdoutPipe()
	_ = cmd.Start()

	scanner := bufio.NewScanner(cmdOut)

	for scanner.Scan() {
		line := scanner.Text()
		fn(line)
	}

	cmd.Wait()
}

// 辅助函数：判断是否用户取消操作
func IsCanceled(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 130
	}
	return false
}
