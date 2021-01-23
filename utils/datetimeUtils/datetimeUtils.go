package datetimeUtils

import "time"

func FormatDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
