package query_cache_helper

import (
    "fmt"
    "errors"
    "os"
    "syscall"
)

const (
    fifoName = "/tmp/dbchain-fifo-1.pipe"
)

var (
    notificationBufferMap = make(map[string]map[string]uint)
    fifoWriter *os.File
)

func getFifoWriter() (*os.File, error) {
    if fifoWriter == nil {
        f, err := os.OpenFile(fifoName, syscall.O_RDWR|syscall.O_NONBLOCK, 0)
        if err != nil {
            fmt.Printf("\n\nError: Failed to open named pipe %s for writing!\n\n", fifoName)
            return nil, errors.New("Failed to open named pipe")
        }
        fifoWriter = f
    }
    return fifoWriter, nil
}

func NotifyTableExpiration(appCode, tableName string) {
    if appCode == "" {
        writer, err := getFifoWriter()
        if err != nil {
            return        // do nothing if failed to get writer
        }

        fmt.Fprintf(writer, "%s,%s\n", "_", "_")

        for k1 := range notificationBufferMap {
            v1 := notificationBufferMap[k1]
            for k2 := range v1 {
                fmt.Fprintf(writer, "%s,%s\n", k1, k2)        // k1: appCode, k2: tableName
                delete(v1, k2)
            }
            delete(notificationBufferMap, k1)
        }
    } else {
        if v, found := notificationBufferMap[appCode]; found {
            v[tableName] = 1
        } else {
            notificationBufferMap[appCode] = map[string]uint{tableName: 1}
        }
    }
}
