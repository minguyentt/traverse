package utils

import (
	"strconv"
)

func ParseUint16(s string) (uint16, error) {
	base := 10
	bitSize := 16

	parse, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		return 0, nil
	}

	return uint16(parse), nil
}

func Int64ToStr(v int64) string {
	return strconv.FormatInt(v, 10)
}
