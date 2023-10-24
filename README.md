# SnowFlake V0.0.4
 Golang implementation of Twitter’s snowflake algorithm

# speed

```golang
// v0.0.3 befer slow!!
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkFib10-4            1208           1005143 ns/op

// v0.0.4 after  very fast!
cpu: Intel(R) Core(TM) i5-5250U CPU @ 1.60GHz
BenchmarkFib10-4           5000000             243 ns/op
```

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
	snowf, err := snowflake.GetSnowFlake(0, "", "")
	if err != nil {
		fmt.Println(err)
	}
	// output ID
	fmt.Println(snowf.Generate())
}
```

#  Introduction 

temporarily the default is to use 2018-01-01 00:00:00 as the starting calculation time

## snowflake.GetSnowFlake() initialization function

```golang
// @param: node int64. between 1~1024 default:0
// @param: location string. example: "Asia/Shanghai" default:""
// @param: startTime string. example: "2022-10-10 10:00:00" default:""
snowflake.GetSnowFlake(0, "", "") 
// or
snowflake.GetSnowFlake(0, "Asia/shaai", "2022-10-10 10:00:00")
```

### The first parameter is of type int64, ranging from 1 to 1024. By default, 0 can be passed in.

```golang
// @param: node int64. between 1~1024 default:0
```

### The second parameter is of string type, example: "Asia/Shanghai" default:"", which can satisfy time.location()

```golang
// Set the time in the specified time zone, the default is local time
// @param location. region For example: Asia/Shanghai Shanghai, China
// @return error. If the time zone fails to be specified, an error message will be returned, otherwise nil will be returned.
func (snowf *SnowFlake) loadLocation(location string) error {
	localSH, err := time.LoadLocation(location)
	if err != nil {
		return errors.New("Specify time zone: [" + location + "] Failed, please check before using")
	}
	snowf.location = localSH
	return nil
}
```

### The third parameter is of string type, example: "2022-10-10 10:00:00" default:"", this format is enough

```golang
// Set the specified start timestamp
// @param startTime. for example:"2022-10-10 10:00:00"
func (snowf *SnowFlake) setStartTime(startTime string) error {
	timeTmeplate := "2006-01-02 15:04:05"
	// Parse datetime string
	times, err := time.ParseInLocation(timeTmeplate, startTime, snowf.location)
	if err != nil {
		return err
	}
	snowf.originTimestamp = times.UnixMilli()
	return nil
}
```


## snowf.Generate() OUT PUT ID