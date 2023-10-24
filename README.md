# SnowFlake
 Golang implementation of Twitter’s snowflake algorithm

# Instructions:

## Download dependency packages

```golang
go get "github.com/bingjiekang/SnowFlake"
```

## Import dependency packages

```golang
import (
    snowflake "github.com/bingjiekang/SnowFlake"
)   
```

## use

```golang
func main(){
	 // initialization
	 tm := snowflake.GetSnowFlake(0, 0)
	 for i := 0; i < 100; i++ {
        // get id
		fmt.Println(tm.NextId())
	}
}
```

#  Introduction (temporarily the default is to use 2018-01-01 00:00:00 as the starting calculation time)

## snowflake.GetSnowFlake() initialization function

```golang
snowflake.GetSnowFlake(0, 0) // Used to initialize the function. The two parameters range from 0 to 31. By default, 0 can be passed in.
```

## .NextId() is used to get different IDs

```golang
snowf := snowflake.GetSnowFlake(0, 0)
fmt.Println(snowf.NextId()) // output ID
```

## .LoadLocation() loads the time zone, the default is local time

```golang
// .LoadLocation() The incoming parameter is a string, refer to time.LoadLocation() to support time zones
snowf := snowflake.GetSnowFlake(0, 0)
snowf.LoadLocation("Asia/Shanghai") // shanghai location
fmt.Println(snowf.NextId()) // output ID
```
