package utils

import "time"

func GetTime() *time.Time {
	currentTime := time.Now()
	var ptrToTime *time.Time
	ptrToTime = &currentTime
	return ptrToTime
}
