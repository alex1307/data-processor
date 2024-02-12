package utils

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/sirupsen/logrus"
)

const DEFAULT_VEHICLE_FILE_NAME = "vehicle"
const DEFAULT_WORKING_DIR = "Software/release/Rust/scraper/resources/data"

func FileName(dir_name string, file_name string) string {
	var working_dir string
	if dir_name == "" {
		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		working_dir = fmt.Sprintf("%s/%s", currentUser.HomeDir, DEFAULT_WORKING_DIR)
	} else {
		working_dir = dir_name
	}
	_, err := os.Stat(working_dir)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Error("Directory Does not exists: ", working_dir)
			panic(err)
		}
	}
	var working_file_name string
	if file_name == "" {
		working_file_name = fmt.Sprintf("%s/%s-%s.csv", working_dir, DEFAULT_VEHICLE_FILE_NAME, time.Now().Format("2006-01-02"))
	} else {
		working_file_name = fmt.Sprintf("%s/%s", working_dir, file_name)
	}
	_, err = os.Stat(working_file_name)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Error("File name {} not exists.", working_file_name)
			panic(err)
		}
	}
	return working_file_name
}

func ReadFiles(folder string, file_name string, extension string) []string {
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	dates := getDatesInRange("2023-05-15", formattedDate)
	files := []string{}
	for _, date := range dates {
		filename := fmt.Sprintf("%s/%s-%s.%s", folder, file_name, date, extension)
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
