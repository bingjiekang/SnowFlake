// Refer to the implementation of twitter snowflake algorithm golang
// The code is for learning reference only
// @author:Jay
package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	// Default start timestamp 2018-01-01 00:00:00
	// You can customize the start timestamp
	originTimestamp int64 = 1514736000000

	// machineBits is the number of working machine id bits
	// Remember, you have a total 22 bits to share between machineBits-sequenceBits
	machineBits uint8 = 10

	// sequenceBits is the number of digits in the serial number id
	sequenceBits uint8 = 12
)

type SnowFlake struct {
	lock            sync.Mutex // Lock
	originTimestamp int64      // origin timestamp
	epoch           time.Time  // current timestamp
	lastTimestamp   int64      // The timestamp of the latest ID
	node            int64
	sequence        int64          // Sequence within milliseconds (0~4095)
	nodeMax         int64          // Maximum machine identification ID
	nodeMask        int64          // machine and serial mask
	sequenceMask    int64          // serial number mask
	timeShift       uint8          // machine and serial Shift left bits
	sequenceShift   uint8          // serial number Shift left bits
	location        *time.Location // Location
}

// Default:(0,"","")
// @param: node int64. between 1~1024 default:0
// @param: location string. example: "Asia/Shanghai" default:""
// @param: startTime string. example: "2022-10-10 10:00:00" default:""
// @return: a new snowflake node that can be used to generate snowflake
func GetSnowFlake(node int64, location string, startTime string) (*SnowFlake, error) {
	// Work machine id and serial number
	if machineBits+sequenceBits > 22 {
		return nil, errors.New("Remember, you have a total 22 bits to share between machineBits-sequenceBits")
	}

	snowf := SnowFlake{}
	snowf.node = node
	snowf.nodeMax = -1 ^ (-1 << machineBits)       // max 1023
	snowf.nodeMask = snowf.nodeMax << sequenceBits // Shift left 22 bits
	snowf.sequenceMask = -1 ^ (-1 << sequenceBits) // 4095
	snowf.timeShift = machineBits + sequenceBits   // machine and serial bits
	snowf.sequenceShift = sequenceBits             // serial number bits

	// node needs to be between 0 and 1024(1<<10) Work machine id bits
	if snowf.node < 0 || snowf.node > snowf.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(snowf.nodeMax, 10))
	}
	// The location is empty by default. If it is not empty, the corresponding time zone is specified.
	if location == "" {
		snowf.loadLocation("")
	} else {
		if err := snowf.loadLocation(location); err != nil {
			return nil, err
		}
	}
	// startTime is empty by default. If it is not empty, the initial time node is specified.
	if startTime == "" {
		snowf.originTimestamp = originTimestamp // start time
	} else {
		if err := snowf.setStartTime(startTime); err != nil {
			return nil, err
		}
	}

	var curTime = time.Now().In(snowf.location)
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	snowf.epoch = curTime.Add(time.Unix(snowf.originTimestamp/1000, (snowf.originTimestamp%1000)*1000000).Sub(curTime))

	return &snowf, nil
}

// Generate creates and returns a unique snowflake ID
// To help guarantee uniqueness
// Make sure your system is keeping accurate system time
// Make sure you never have multiple nodes running with the same node ID
// @return snowflakeId
func (snowf *SnowFlake) Generate() int64 {

	snowf.lock.Lock()
	defer snowf.lock.Unlock()
	// current time
	timestamp := time.Since(snowf.epoch).Milliseconds()
	// If they are generated at the same time, sort the sequence within milliseconds.
	if timestamp == snowf.lastTimestamp {
		snowf.sequence = (snowf.sequence + 1) & snowf.sequenceMask
		// Sequence overflow exceeds 4095 in milliseconds
		if snowf.sequence == 0 {
			for timestamp <= snowf.lastTimestamp {
				timestamp = time.Since(snowf.epoch).Milliseconds()
			}
		}
	} else {
		// Timestamp changes, timestamp resets within milliseconds
		snowf.sequence = 0
	}
	// The last time the ID was generated
	snowf.lastTimestamp = timestamp

	// Shift and pass or run to spell the 64-bit ID
	return int64((timestamp)<<snowf.timeShift |
		(snowf.node << snowf.sequenceShift) |
		(snowf.sequence),
	)
}

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
