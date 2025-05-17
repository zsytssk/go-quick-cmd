package command

import (
	"fmt"
	"io"
	"log"
	"quick-cmd/dbt"
	"quick-cmd/utils"
	"strings"
)

func JumpDir() {
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	dbPath, err := utils.GetCurDirFileName("db")
	if err != nil {
		log.Fatal(err)
	}
	db, err := dbt.Init(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	items, err := dbt.GetDir(db)
	if err != nil {
		log.Fatal(err)
	}

	cmdStr := buildFindStr(config)
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		for _, item := range items {
			fmt.Fprintf(writer, "%s [%d:%d]\n", item.Name, item.ID, item.Priority)
		}
		utils.RunCMDInSteam(cmdStr, func(line string) {
			index := utils.ArrFindIndex(items, func(item dbt.Item, _ int) bool {
				return item.Name == line
			})
			if index != -1 {
				return
			}
			item := dbt.Item{ID: -1, Name: line, Priority: 0}
			items = append(items, item)
			fmt.Fprintf(writer, "%s [%d:%d]\n", item.Name, item.ID, item.Priority)
		})
	}()

	selected, err := utils.RunFZFStream(reader)

	if err != nil {
		if utils.IsCanceled(err) { // 检查是否用户取消
			log.Println("选择已取消")
			return
		}
		log.Fatal(err)
	}

	if selected == "" {
		log.Println("未选择任何项目")
		return
	}

	index := utils.ArrFindIndex(items, func(item dbt.Item, _ int) bool {
		return selected == fmt.Sprintf("%s [%d:%d]", item.Name, item.ID, item.Priority)
	})

	if index == -1 {
		log.Println("find item index = ", selected)
		return
	}
	item := items[index]
	dbt.UpdateDirPriority(db, item)

	fmt.Print(`cd `, item.Name)
}

func buildFindStr(config utils.Config) string {
	var cmdInput strings.Builder
	for _, item := range config.Folders {
		ignores := append(config.Ignores, item.Ignores...)
		ignoreStr := utils.ArrJoin(ignores, func(item string, index int) string {
			if index == 0 {
				return fmt.Sprintf(` -path %s`, item)
			}
			return fmt.Sprintf(` -o -path %s`, item)
		})
		cmdInput.WriteString(fmt.Sprintf("find %s -maxdepth %d -type d \\( %s \\)  -prune -o -print\n", item.Folder, item.Depth, ignoreStr))
	}

	return cmdInput.String()
}
