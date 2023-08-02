package utils

import "time"

func ConvertDate(value string) time.Time {

	if value == "" {
		return time.Now()
	}

	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Now()
	}

	return date
}
