package main

import (
	"data-processor/app"
	"data-processor/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"runtime/pprof"
	"time"
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
	defaultDir := fmt.Sprintf("%s/%s", currentUser.HomeDir, utils.DEFAULT_WORKING_DIR)
	defaultFileName := fmt.Sprintf("%s-%s.csv", utils.DEFAULT_VEHICLE_FILE_NAME, time.Now().Format("2006-01-02"))
	defaultUpdatedName := fmt.Sprintf("%s-%s.csv", "updated-vehicle", time.Now().Format("2006-01-02"))
	log.Println("Default directory: ", defaultDir)
	update := flag.NewFlagSet("update", flag.ExitOnError)
	updateFileName := update.String("source", defaultUpdatedName, "Vehicle csv source file name")
	updatedDataDirName := update.String("dir", defaultDir, "File directory name")

	add := flag.NewFlagSet("add", flag.ExitOnError)
	addFileName := add.String("source", defaultFileName, "Vehicle csv source file name")
	metaDataFileName := add.String("meta-data", "meta_data.csv", "Search meta-data source file name")
	newDataDirName := add.String("dir", defaultDir, "File directory name")

	delete := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteFileName := delete.String("source", "not-found-ids.csv", "Vehicle csv source file name")
	deletedFilesDirName := delete.String("dir", defaultDir, "File directory name")

	forUpdate := flag.NewFlagSet("update-list", flag.ExitOnError)
	forUpdateFileName := forUpdate.String("source", "ids-for-update.csv", "Vehicle csv source file name")
	forUpdateFilesDirName := forUpdate.String("dir", defaultDir, "File directory name")

	if len(os.Args) < 2 {

	} else if os.Args[1] == "update" {
		update.Parse(os.Args[2:])
		log.Println("Updated vehicles directory: ", *updatedDataDirName)
		log.Println("Updated vehicles file name: ", *updateFileName)
		app.UpdateVehicles(*updatedDataDirName, *updateFileName)
	} else if os.Args[1] == "add" {
		add.Parse(os.Args[2:])
		log.Println("New vehicles file name: ", *addFileName)
		log.Println("Meta data file name: ", *metaDataFileName)
		log.Println("New vehicles directory: ", *newDataDirName)
		app.AddNewVehicles(*newDataDirName, *metaDataFileName, *addFileName)
	} else if os.Args[1] == "delete" {
		delete.Parse(os.Args[2:])
		log.Println("Deleted vehicles file name: ", *deleteFileName)
		log.Println("Deleted vehicles directory: ", *deletedFilesDirName)
		app.DeletedVehiclesByIds(*deletedFilesDirName, *deleteFileName)
	} else if os.Args[1] == "update-list" {
		forUpdate.Parse(os.Args[2:])
		log.Println("Vehicles for update file name: ", *forUpdateFileName)
		log.Println("Vehicles for update directory: ", *forUpdateFilesDirName)
		app.GenerateUpdateList(*forUpdateFilesDirName, *forUpdateFileName)
	} else {
		log.Println("Default: Process new vehicles... ", defaultDir)
		var dirName string
		var fileName string
		var metaDataFileName string

		flag.StringVar(&fileName, "source", defaultFileName, "Vehicle csv source file name")
		flag.StringVar(&metaDataFileName, "meta-data", "meta_data.csv", "Meta data csv source file name")
		flag.StringVar(&dirName, "dir", defaultDir, "Data directory")
		flag.Parse()
		app.AddNewVehicles(dirName, metaDataFileName, fileName)
	}
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
