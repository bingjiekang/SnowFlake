# SnowFlake
 Twitter雪花算法的Golang实现

# 使用方法:

```golang
import (
    "github.com/bingjiekang/SnowFlake"
)
    // 初始化
    tm := SnowFlake.GetSnowFlake(0, 0)
	for i := 0; i < 100; i++ {
        // 获得id
		fmt.Println(tm.NextId())
	}

```
