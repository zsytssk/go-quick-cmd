package command

import (
	"fmt"
	"log"
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
	var cfg *ini.File
	var topPath string
	var err error
	var localErr error
	for i := 0; i < 2; i++ {
		cmd := "git rev-parse --show-toplevel"
		if i == 1 {
			cmd = "cd .. && git rev-parse --show-toplevel"
		}
		topPath, localErr = utils.RunCMD(cmd)
		if localErr != nil && i == 0 {
			log.Fatalf("无法读取 git目录: %v", localErr)
		}

		cfg, localErr = ini.Load(fmt.Sprintf(`%s/.gitmodules`, topPath))
		if cfg != nil {
			break
		}
		if i == 0 {
			err = localErr
		}
	}

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
		Items: list,
		// Path:  "/home/zsy/Documents/zsy/job/plims-background",
		Path: topPath,
	}
}

func runSubmodulesCmd(submodule Submodules, cmd string, lite bool) (err error) {
	if !lite {
		fmt.Println(submodule.Path)
	}
	output, err := utils.RunCMD(fmt.Sprintf("cd %s && %s", submodule.Path, cmd))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

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
