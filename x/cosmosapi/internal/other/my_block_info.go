package other

import (
    "time"
)

var (
    blockHeight int64
    blockTime time.Time
)

func SaveCurrentBlockInfo(height int64, timeStamp time.Time) {
    blockHeight = height
    blockTime = timeStamp
}

func GetcurrentBlockHeight() int64 {
    return blockHeight
}

func GetcurrentBlockTime() time.Time {
    return blockTime
}
