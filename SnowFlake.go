// Refer to the implementation of twitter snowflake algorithm golang
// The code is for learning reference only
// @author:Jay
package SnowFlake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type SnowFlake struct {
	// 机器ID(0~31)
	machineId int64
	// 数据中心ID(0~31)
	dataCenterId int64
	// 地区
	location *time.Location
}

var (
	// 默认开始时间戳 2018-01-01 00:00:00
	originTimestamp int64 = 1514736000000
	// 数据中心ID所占位数 5位
	dataCenterIdBits int64 = 5
	// 机器标识所占位数 5位
	machineWorkerIdBits int64 = 5
	// 最大数据中心ID 31
	maxDataCenterId int64 = -1 * (-1 << dataCenterIdBits)
	// 最大机器标识ID 31
	maxMachineWorkerId int64 = -1 * (-1 << machineWorkerIdBits)
	// 序列号所占位数 12
	sequenceBits int64 = 12
	// 机器ID左移12位(序列号位数)
	MachineIdShift = sequenceBits
	// 数据中心ID左移17位(序列号+机器ID位)
	datacenterIdShift = sequenceBits + machineWorkerIdBits
	// 时间戳向左移动22位(序列号+机器ID位+数据中心ID位)
	timestampShift = sequenceBits + machineWorkerIdBits + dataCenterIdBits
	// 生成序列掩码 4095<-(0b111111111111=0xfff)<-(1<<12)
	sequenceMask int64 = -1 * (-1 << sequenceBits)
	// 机器ID(0~31)
	machineId int64
	// 数据中心ID(0~31)
	dataCenterId int64
	// 毫秒内序列(0~4095)
	sequence int64 = 0
	// 上次生成ID的时间戳
	lastTimestamp int64 = -1
)

// 构造函数 默认都传入0即可
// @param: MachineId 机器ID(0~31)
// @param: dataCenterId 数据中心ID(0~31)
func GetSnowFlake(machineId, dataCenterId int64) SnowFlake {
	var snowf SnowFlake
	if machineId > maxMachineWorkerId || machineId < 0 {
		panic(fmt.Sprintf("machine Id can't be greater than %d or less than 0", maxMachineWorkerId))
	}
	if dataCenterId > maxDataCenterId || dataCenterId < 0 {
		panic(fmt.Sprintf("dataCenter Id can't be greater than %d or less than 0", maxDataCenterId))
	}
	snowf.machineId = machineId
	snowf.dataCenterId = dataCenterId
	return snowf
}

// 获得下一个ID
// @return snowflakeId
func (snowf *SnowFlake) NextId() int64 {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	timestamp := snowf.timeNow()
	// 如果当前时间小于上一次ID生成的时间戳,说明系统时钟回退过,抛出异常Panic
	if timestamp < lastTimestamp {
		panic(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", lastTimestamp-timestamp))
	}
	// 如果是同一时间生成,进行毫秒内序列sequence排序
	if timestamp == lastTimestamp {
		sequence = (sequence + 1) & sequenceMask
		// 毫秒内序列溢出 超过4095个
		if sequence == 0 {
			// 阻塞到下一个毫秒,获得新的时间戳
			timestamp = snowf.timeNextMilli(lastTimestamp)
		}
	} else { // 时间戳改变,毫秒内时间戳重置
		sequence = 0
	}
	// 上次生成ID的时间截
	lastTimestamp = timestamp
	// 移位并通过 或运行 拼成64位ID
	return ((timestamp-originTimestamp)<<timestampShift | (dataCenterId << datacenterIdShift) | (machineId << MachineIdShift) | sequence)

}

// 获取以毫秒为单位的当前时间
// @param timeZone 时区localSH, _ := time.LoadLocation("Asia/Shanghai")
// @return 当前毫秒时间
func (snowf *SnowFlake) timeNow() int64 {
	if snowf.location == nil {
		return time.Now().UnixMilli()
	}
	return time.Now().In(snowf.location).UnixMilli()
}

// 阻塞到下一个毫秒，直到获得新的时间戳
// @param lastTimestamp 上次生成ID的时间截
// @return 当前时间戳
func (snowf *SnowFlake) timeNextMilli(lastTimestamp int64) int64 {
	timestamp := snowf.timeNow()
	for timestamp <= lastTimestamp {
		timestamp = snowf.timeNow()
	}
	return timestamp
}

// 设置指定时区的时间,默认为本地时间
// @param location 地区 例如:Asia/Shanghai 中国上海
// @return error 指定时区失败返回错误信息,否则返回nil
func (snowf *SnowFlake) LoadLocation(location string) error {
	localSH, err := time.LoadLocation(location)
	if err != nil {
		return errors.New("指定时区:" + location + "失败,请检查后使用")
	}
	snowf.location = localSH
	return nil
}
