package utils

import (
    "unsafe"
)

func IntToByteArray(num int64) []byte {
    size := int(unsafe.Sizeof(num))
    arr := make([]byte, size)
    for i := 0 ; i < size ; i++ {
        byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
        arr[size-1-i] = byt
    }
    return arr
}

func ByteArrayToInt(arr []byte) int64 {
    val := int64(0)
    size := len(arr)
    for i := 0 ; i < size ; i++ {
        *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[size-1-i]
    }
    return val
}

