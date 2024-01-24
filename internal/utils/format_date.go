package utils

import (
	"fmt"
	"time"
)

func FormatDate(date time.Time) string {
	year, month, day := date.Date()
	monthStr := fmt.Sprintf("%02d", int(month))
	dayStr := fmt.Sprintf("%02d", day)

	return fmt.Sprintf("%d-%s-%s", year, monthStr, dayStr)
}
