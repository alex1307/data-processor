package utils

import (
	"fmt"
	"os"
	"time"
)

func ReadFiles(folder string, file_name string, extension string) []string {
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	dates := getDatesInRange("2023-05-15", formattedDate)
	files := []string{}
	for _, date := range dates {
		filename := fmt.Sprintf("%s/%s_%s.%s", folder, file_name, date, extension)
		_, err := os.Stat(filename)
		if !os.IsNotExist(err) {
			files = append(files, filename)
		}
	}
	return files
}

func getDatesInRange(date1, date2 string) []string {
	layout := "2006-01-02"

	start, _ := time.Parse(layout, date1)
	end, _ := time.Parse(layout, date2)

	dates := []string{}

	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		dates = append(dates, current.Format(layout))
	}

	return dates
}
