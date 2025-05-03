- @todo 读取 history 写入数据库

## 2025-05-02 17:51:49

- @ques 也许我可以写一个脚本 拉取仓库 然后编辑 然后导出执行文件到某个地方

  - 然后再把那些执行文件 ignore 了，这样就不用担心编译后文件太大的问题了

- @ques 编译到 awesome/bin

- @todo
  - 支持多个命令

### end

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
