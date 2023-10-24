# SnowFlake
 Twitter雪花算法的Golang实现

# 使用方法:

## 下载依赖包

```golang
go get "github.com/bingjiekang/SnowFlake"
```

## 导入依赖包

```golang
import (
    "github.com/bingjiekang/SnowFlake"
)   
```

## 使用

```golang
func main(){
	 // 初始化
	 tm := SnowFlake.GetSnowFlake(0, 0)
	 for i := 0; i < 100; i++ {
        // 获得id
		fmt.Println(tm.NextId())
	}
}
```

#  介绍

## SnowFlake.GetSnowFlake() 初始化函数

```golang
SnowFlake.GetSnowFlake(0, 0) // 用来初始化函数，两个参数范围在0~31，默认传入0即可
```

## .NextId() 用来获取不同ID

```golang
snowf := SnowFlake.GetSnowFlake(0, 0)
fmt.Println(snowf.NextId()) // output ID
```

## .LoadLocation() 加载时区，默认为本机时间

```golang
snowf := SnowFlake.GetSnowFlake(0, 0)
snowf.LoadLocation("Asia/Shanghai") // shanghai location
fmt.Println(snowf.NextId()) // output ID

// .LoadLocation() 传入参数为字符串，参考time.LoadLocation()支持时区
```
