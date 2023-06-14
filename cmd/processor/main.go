package main

import (
	"data-processor/app"
	"fmt"
	"log"
	"os/user"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	DataDir := fmt.Sprintf("%s/%s", currentUser.HomeDir, "Software/release/Rust/scraper/resources/data")
	log.Println("Working data directory: ", DataDir)
	app.ProcessCSVFiles(DataDir)
}
