package util

import "time"

func Now() time.Time {
	return time.Now()
}

func NowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}
