package utils

import (
	"strconv"
	"time"
)

func MsToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.UnixMilli(msInt), nil
}
