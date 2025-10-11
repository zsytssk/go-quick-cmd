package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"quick-cmd/utils"
	"strings"

	"gopkg.in/ini.v1"
)

var supportGitCmd = []string{"submodule", "submoduleLite"}

func Git(cmd []string) (err error) {
	if len(cmd) == 0 {
		fmt.Println(`请输入git下级命令`)
		return
	}

	switch cmd[0] {
	case "submodule":
		if len(cmd) == 1 {
			fmt.Println(`请输入git submodule下级命令`)
			return
		}
		submodules := getSubmodules()
		err = runSubmodulesCmd(submodules, strings.Join(cmd[1:], " "), false)
	case "submoduleLite":
		if len(cmd) == 1 {
			fmt.Println(`请输入git submodule下级命令`)
			return
		}
		submodules := getSubmodules()
		err = runSubmodulesCmd(submodules, strings.Join(cmd[1:], " "), true)
	default:
		fmt.Println(`只支持命令:`, strings.Join(supportGitCmd, ", "))
		return
	}

	return
}

type Submodules struct {
	IsGit bool
	Path  string
	Items []SubmoduleItem
}
type SubmoduleItem struct {
	Name   string
	Path   string
	URL    string
	Branch string
}

func getSubmodules() Submodules {

	cfg, topPath, err := findGitModules()
	if cfg == nil {
		log.Fatalf("无法读取 .gitmodules: %v", err)
	}
	var list []SubmoduleItem
	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}
		list = append(list, SubmoduleItem{
			section.Name(),
			section.Key("path").String(),
			section.Key("url").String(),
			section.Key("branch").String(),
		})
	}
	return Submodules{
		IsGit: isGitRepo(topPath),
		Items: list,
		Path:  topPath,
	}
}

func runSubmodulesCmd(submodule Submodules, cmd string, lite bool) (err error) {
	if !lite {
		fmt.Println(submodule.Path)
	}
	if submodule.IsGit {
		output, err := utils.RunCMD(fmt.Sprintf("cd %s && %s", submodule.Path, cmd))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	}

	for _, item := range submodule.Items {
		fullPath := fmt.Sprintf("%s/%s", submodule.Path, item.Path)
		if !lite {
			fmt.Printf("---\n%s\n", fullPath)
		}
		output, err := utils.RunCMD(fmt.Sprintf("cd %s && %s", fullPath, cmd))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	}
	return
}

func findGitModules() (*ini.File, string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, "", fmt.Errorf("获取当前工作目录失败: %w", err)
	}
	checkFileNames := []string{".gitmodules.local", ".gitmodules"}
	for {
		// 检查两个可能的文件
		for _, name := range checkFileNames {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				cfg, cfgErr := ini.Load(path)
				return cfg, dir, cfgErr
			}
		}

		// 上移目录
		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到达根目录
			break
		}
		dir = parent
	}

	return nil, "", fmt.Errorf("未找到 .gitmodules.local 或 .gitmodules 文件")
}

func isGitRepo(dir string) bool {
	gitPath := filepath.Join(dir, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return false
	}
	// .git 可以是目录（普通仓库）或文件（子模块等）
	return info.IsDir() || !info.IsDir()
}
