package datetime

import (
	"time"
)

func StrParseToTime(dateValue, dateLayout string) time.Time {
	dateResult, _ := time.Parse(dateLayout, dateValue)

	return dateResult
}

func EmptyTime(timeInput time.Time) *time.Time{
	if timeInput.IsZero() {
		return nil
	}

	return &timeInput
}
