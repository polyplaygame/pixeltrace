package tools

import (
	"fmt"
	"time"
)

func TimezoneAndHourOffset(locationName string) (string, float64, error) {
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		return "", 0, fmt.Errorf("failed to load location: %v", err)
	}
	tz, offset := time.Now().In(loc).Zone()
	// offset is in seconds, convert to hours，要求返回小时，不满一小时的处理为小数
	return tz, float64(offset) / 3600, nil
}
