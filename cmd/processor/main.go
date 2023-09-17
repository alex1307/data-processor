package main

import (
	"data-processor/app"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"runtime/pprof"
)

func main() {

	runCPUProfile()
	defer runMemProfile()

	runMemProfile()
	defer runMemProfile()

	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	DataDir := fmt.Sprintf("%s/%s", currentUser.HomeDir, "Software/release/Rust/scraper/resources/data")
	log.Println("Working data directory: ", DataDir)
	app.ProcessCSVFiles(DataDir)
}

func runMemProfile() {
	f, err := os.Create("memprofile")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()
}

func runCPUProfile() {
	f, err := os.Create("cpuprofile")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
}
