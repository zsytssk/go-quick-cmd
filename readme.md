- @todo 读取 history 写入数据库

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
