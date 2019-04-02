package helper

import (
	"strconv"
)

func Int2str(bs []uint8) string {
	return string(bs)
}

func Int2int64(bs []uint8) int64 {
	str := Int2str(bs)
	r, _ := strconv.ParseInt(str, 10, 64)
	return r
}
