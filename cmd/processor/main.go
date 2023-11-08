package main

import (
	"data-processor/app"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {

	runCPUProfile()
	defer runMemProfile()

	runMemProfile()
	defer runMemProfile()

	update := flag.NewFlagSet("update", flag.ExitOnError)
	updateFileName := update.String("source", "", "Vehicle csv source file name")
	updatedDataDirName := update.String("dir", "", "File directory name")

	add := flag.NewFlagSet("add", flag.ExitOnError)
	addFileName := add.String("source", "", "Vehicle csv source file name")
	metaDataFileName := add.String("meta-data", "meta_data.csv", "Search meta-data source file name")
	newDataDirName := add.String("dir", "", "File directory name")

	delete := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteFileName := delete.String("source", "", "Vehicle csv source file name")
	deletedFilesDirName := delete.String("dir", "", "File directory name")
	if len(os.Args) < 2 {

	} else if os.Args[1] == "update" {
		update.Parse(os.Args[2:])
		log.Println("Updated vehicles file name: ", *updateFileName)
		log.Println("Updated vehicles directory: ", *updatedDataDirName)
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
	} else {
		var dir string
		var vehicles_file_name string
		var meta_data_file_name string
		var for_update_file_name string
		var deleted_file_name string

		update := flag.NewFlagSet("update", flag.ExitOnError)
		update.String("source", "", "Vehicle csv source file name")

		add := flag.NewFlagSet("add", flag.ExitOnError)
		add.String("source", "", "Vehicle csv source file name")

		delete := flag.NewFlagSet("delete", flag.ExitOnError)
		delete.String("source", "", "Vehicle csv source file name")

		flag.StringVar(&vehicles_file_name, "source", "", "Vehicle csv source file name")
		flag.StringVar(&meta_data_file_name, "meta-data-file", "", "Meta data csv source file name")
		flag.StringVar(&for_update_file_name, "for-update-file", "", "Adverts for update file name")
		flag.StringVar(&deleted_file_name, "deleted-ids-file", "", "Deleted adverts csv file name")
		flag.StringVar(&dir, "data_dir", "", "Data directory")
		flag.Parse()
		log.Println("Data directory: ", dir)
		log.Println("Vehicles file name: ", vehicles_file_name)
		log.Println("Meta data file name: ", meta_data_file_name)
		log.Println("For update file name: ", for_update_file_name)
		log.Println("Deleted ids file name: ", deleted_file_name)
	}

	// metasearh_file_name := fmt.Sprintf("%s/%s", DataDir, "meta_data.csv")
	// records_for_update_file_name := fmt.Sprintf("%s/%s", DataDir, "for_update.csv")
	// updated_records_file_name := fmt.Sprintf("%s/%s", DataDir, "updated.csv")
	// deleted_file_name := fmt.Sprintf("%s/%s", DataDir, "deleted.csv")

	// log.Println("Working data directory: ", DataDir)
	// log.Println("args: {}, len: {}", os.Args, len(os.Args))
	// // if len(os.Args) < 2 {
	// // 	log.Println("No file name provided. Processing all files in the data directory")
	// // 	app.ProcessCSVFiles(DataDir, "")
	// // } else {
	// // 	app.ProcessDeletedIds(DataDir, os.Args[1])
	// // }
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
