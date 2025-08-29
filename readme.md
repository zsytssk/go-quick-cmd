`go build  -o /home/zsy/.config/awesome/autostart/ignoreBin/quick-cmd`

- @ques 有什么像 fzf 但是能设置权重的工具

## 2025-07-14 14:44:03

- action 能不能写成类型

---

- @todo hide
  - fzf 监听快捷键

```
lines := strings.Split(buf.String(), "\n")
key := lines[0] // ctrl-d
choice := lines[1] // Cherry
```

- @todo 集合命令 -> 通过名称区分
  - unionCommand
  - name content
  - 激活到 bashHistory 中

```
ghostty --working-directory=/home/zsy/Documents/zsy/job/wms-background/web -e "npm run serve" & ghostty --working-directory=/home/zsy/Documents/zsy/job/wms-background/server -e "go build main.go && ./main"
```

```
nohup ghostty --working-directory=/home/zsy/Documents/zsy/job/wms-background/web \
  -e "npm run serve" > /dev/null 2>&1 &

nohup ghostty --working-directory=/home/zsy/Documents/zsy/job/wms-background/server \
  -e "go build main.go && ./main" > /dev/null 2>&1 &
```

## 2025-07-13 21:19:49

---

- @todo struct -> map
  - 支持 tag json
- @todo map -> struct

`FieldItem` 实际应用
`StructToSQLCreateTable` -> 默认值 `NOT NULL DEFAULT`
IsSQLTypeCompatible 改写
SyncTableColumns 只要执行一次

## 2025-05-02 17:51:49

- @ques 能不能像后端代码一样通过定义 struct 来控制 table 和插入更新数据？

```
InitTable
GetList
Delete
Update
Insert
```

```go
func (DeviceShutdownOperation) TableName() string {
	return "device_shutdown_operation"
}
```

### end

- @ques 为了避免卡顿 能不能使用 stream 传递数据?

  - 可行
  - 怎么把数据库的放在最底下

- @ques 扩展 update 方法 支持任意数量属性
- @ques 修改包名字

- @ques 一键编译...

  - 编译到 awesome/bin

- @todo 添加全局快捷命令

  - 保存 init
  - 编译脚本等等

- @todo 支持快捷键删除

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
