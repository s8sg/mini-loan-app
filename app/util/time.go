package util

import "time"

func GetCurrentTimeInUtc() time.Time {
	return time.Now().UTC()
}
