`go build  -o /home/zsy/.config/awesome/autostart/ignoreBin/quickCmd`

- @ques 有什么像 fzf 但是能设置权重的工具

## 2025-05-02 17:51:49

- @ques 修改包名字

- @ques 一键编译...

  - 编译到 awesome/bin

- @todo 支持快捷键删除

- @ques 为了避免卡顿 能不能使用 stream 传递数据?

- @ques 扩展 update 方法 支持任意数量属性

### end

```
DROP TABLE table_name;

```

- @todo

  - 支持跳转文件夹命令
  - 记录跳转文件夹

- 支持记录 文件夹

  - 历史记录中的常见的

- bashHistory 排除 `cd`

```bash
(
  find dirA -maxdepth 1 ! -name "*.log"
  find dirB -maxdepth 3 -path "*/tmp" -prune -o -print
  find dirC -path "*/.git" -prune -o -print
) | fzf
```

```bash
(
  find ~/.config -maxdepth 1 -path "*node_modules*"  -path "*.git*"  -prune -o -print
  find ~/Documents/zsy/ -maxdepth 2 -path "*node_modules*"  -path "*.git*" -prune -o -print
) | fzf
```

- @todo 读取 history 写入数据库
- @ques 插入数据
- @ques 也许我可以写一个脚本 拉取仓库 然后编辑 然后导出执行文件到某个地方

  - 然后再把那些执行文件 ignore 了，这样就不用担心编译后文件太大的问题了

- 查找所有命令

- @ques 匹配字符
- @ques fzf 能不能完全匹配按照顺序排列
- @ques 更新 priority
- @todo 检查 table 是否存在 name, 覆盖数据

## 2025-05-02 17:21:11

go 有没有内置像 fzf 功能的包

```go
import "github.com/ktr0731/go-fuzzyfinder"

func main() {
    items := []string{"apple", "banana", "cherry"}

    // 单选模式
    idx, _ := fuzzyfinder.Find(items, func(i int) string {
        return items[i]
    })

    fmt.Printf("Selected: %s\n", items[idx])
}
```
